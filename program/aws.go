package program

import (
	"bufio"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"sync"
)

type sessionInfo struct {
	profile string
	account string
	region  string
	session *session.Session
	log     zerolog.Logger
}

type ClusterInfo struct {
	*eks.Cluster
	log     zerolog.Logger
	session *sessionInfo
}

// getProfiles gets all profiles from ~/.aws/credentials or the program arguement
func (program *Options) getProfiles() <-chan string {
	output := make(chan string)

	if len(program.Profiles) < 1 {
		go func() {
			defer close(output)
			if f, err := os.Open(program.CredentialsFile); err == nil {
				scanner := bufio.NewScanner(f)
				scanner.Split(bufio.ScanLines)

				for scanner.Scan() {
					s := strings.TrimSpace(scanner.Text())
					if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
						output <- s[1 : len(s)-1]
					}
				}

			} else {
				stats.Errors.Add(1)
				log.Error().Str("file", program.CredentialsFile).Err(err).Msg("Failed to open file")
			}
		}()
	} else {
		go func() {
			defer close(output)
			for _, p := range program.Profiles {
				output <- p
			}
		}()
	}

	return output
}

// getClustersFrom gets the clusters from
func (program *Options) getClustersFrom(s *sessionInfo, clusters chan<- ClusterInfo) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	e := eks.New(s.session)

	s.log.Debug().Msg("Listing EKS clusters")
	if out, err := e.ListClusters(&eks.ListClustersInput{}); err == nil {
		s.log.Debug().Msg("Getting Clusters")

		stats.Clusters.Add(int32(len(out.Clusters)))

		for _, c := range out.Clusters {
			wg.Add(1)
			go func(c *string) {
				defer wg.Done()
				s.log.Debug().Str("cluster_name", *c).Msg("Found cluster")

				if out, err := e.DescribeCluster(&eks.DescribeClusterInput{Name: c}); err != nil {
					stats.Errors.Add(1)
					log.Error().Err(err).Msg("Error describing cluster")
				} else {
					log.Info().Str("cluster_name", *c).Str("Profile", s.profile).Str("Region", s.region).Str("Account", s.account).Msg("Cluster config downloaded for")
					clusters <- ClusterInfo{
						Cluster: out.Cluster,
						log:     s.log.With().Str("cluster_name", *c).Logger(),
						session: s,
					}
				}
			}(c)
		}
	} else {
		stats.Errors.Add(1)
		s.log.Error().Err(err).Msg("Error listing clusters")
	}
}

func (program *Options) getUniqueSessions() <-chan *sessionInfo {

	sessions := make(chan *sessionInfo)

	go func() {
		wg := sync.WaitGroup{}

		defer close(sessions)
		defer wg.Wait()

		stats.Regions.Add(int32(len(program.Regions)))

		accounts := make(map[string]string)
		for info := range program.getProfileSessions() {
			if _, found := accounts[info.account]; found {
				info.log.Debug().Msg("Profile is duplicate")
			} else {
				info.log.Debug().Msg("Profile is good for use")
				accounts[info.account] = info.profile

				stats.UniqueProfiles.Add(1)
				sessions <- info

				for _, region := range program.Regions {
					wg.Add(1)
					go func(profile, region, account string) {
						defer wg.Done()
						if region != info.region {
							log := log.With().Str("profile", info.profile).Str("region", region).Logger()
							log.Debug().Msg("Creating regional session")
							if s, err := session.NewSessionWithOptions(session.Options{Profile: profile, Config: aws.Config{Region: aws.String(region)}}); err == nil {
								sessions <- &sessionInfo{
									profile: profile,
									region:  region,
									account: account,
									session: s,
									log:     log,
								}
							} else {
								stats.Errors.Add(1)
								log.Error().Err(err).Msg("Failed to create session")
							}
						}
					}(info.profile, region, info.account)
				}
			}
		}
	}()

	return sessions
}

// getProfileSessions gets a channel for a session for the first region for each profile, and fills in the account ID for that profile.
func (program *Options) getProfileSessions() <-chan *sessionInfo {

	sessions := make(chan *sessionInfo)
	wg := sync.WaitGroup{}

	go func() {
		defer close(sessions)
		defer wg.Wait()

		profiles := program.getProfiles()

		for p := range profiles {
			log := log.With().Str("profile", p).Str("region", program.Regions[0]).Logger()
			stats.Profiles.Add(1)
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				if s, err := NewSession(p, program.Regions[0], log); err == nil {
					stats.UsableProfiles.Add(1)
					sessions <- s
				}
			}(p)

		}
	}()

	return sessions
}

func NewSession(profile, region string, log zerolog.Logger) (*sessionInfo, error) {
	if sess, err := session.NewSessionWithOptions(session.Options{Profile: profile, Config: aws.Config{Region: aws.String(region)}}); err == nil {
		svc := sts.New(sess)
		if out, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{}); err == nil {
			log := log.With().Str("account", *out.Account).Logger()
			log.Debug().Msg("Profile for account")
			return &sessionInfo{
				profile: profile,
				region:  region,
				account: *out.Account,
				session: sess,
				log:     log,
			}, nil
		} else {
			log.Error().Err(err).Msg("Error reaching AWS")
			return nil, err
		}
	} else {
		return nil, err
	}
}
