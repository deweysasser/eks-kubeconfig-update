package program

import (
	"github.com/rs/zerolog/log"
	"sync/atomic"
)

type Stats struct {
	// ProfileRegionPairs is the number of regions checked
	Profiles, UniqueProfiles, UsableProfiles, Regions, Clusters, Errors atomic.Int32
}

func (s *Stats) Log() {
	log.Info().
		Int32("profiles", s.Profiles.Load()).
		Int32("unique_profiles", s.UniqueProfiles.Load()).
		Int32("usable_profiles", s.UsableProfiles.Load()).
		Int32("regions", s.Regions.Load()).
		Int32("clusters", s.Clusters.Load()).
		Int32("fatal_errors", s.Errors.Load()).
		Msg("Statistics")
}

var stats Stats
