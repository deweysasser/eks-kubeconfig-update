# eks-kubeconfig-update

A quick utility for finding and downloading kubeconfig files for all reachable EKS clusters.

## Overview

This utility will download a `kubeconfig` file for all EKS clusters it can locate through all AWS
profiles in your `~/.aws/credentials` file. It does as much as possible in parallel so overall
runtime is very small for the number of profiles and regions retrieved.

It will tolerate redundant profiles.

## Installation

On MacOS: `brew install deweysasser/tap/eks-kubeconfig-update`

(if you happen to use brew on linux, you can also use the above)

On Linux or Windows:  Download the appropriate package from the 
[latest release](https://github.com/deweysasser/eks-kubeconfig-update/releases) page.

## Quick Start

```shell
eks-kubeconfig-update
```

This will use every profile in your `~/.aws/credentials` and download kubeconfig for every cluster located.

## Example

```shell
$ time eks-kubeconfig-update
8:44PM ERR Error reaching AWS error="InvalidClientTokenId: The security token included in the request is invalid\n\tstatus code: 403, request id: e7b25ace-2f73-45a1-9c81-0f0b9e34ba6e" profile=profile-1 region=us-east-1
8:44PM ERR Error reaching AWS error="InvalidClientTokenId: The security token included in the request is invalid\n\tstatus code: 403, request id: 228ece97-cdba-4de2-be70-48135ccad188" profile=profile-2 region=us-east-1
8:44PM ERR Error reaching AWS error="InvalidClientTokenId: The security token included in the request is invalid\n\tstatus code: 403, request id: 6b8e9d80-b024-4db4-8fb4-9918744a03b4" profile=profile-3 region=us-east-1
8:44PM ERR Error reaching AWS error="InvalidClientTokenId: The security token included in the request is invalid.\n\tstatus code: 403, request id: fb285c8f-d819-462e-b1fc-81b7bfde02fe" profile=profile-4 region=us-east-1
8:44PM ERR Error reaching AWS error="InvalidClientTokenId: The security token included in the request is invalid\n\tstatus code: 403, request id: 731fa6c2-43ec-4454-9720-438171faab3c" profile=profile-5 region=us-east-1
8:44PM ERR Error reaching AWS error="InvalidClientTokenId: The security token included in the request is invalid\n\tstatus code: 403, request id: 7f82d496-1d92-466e-bb93-58a48d75a06d" profile=profile-6 region=us-east-1
8:44PM ERR Error reaching AWS error="InvalidClientTokenId: The security token included in the request is invalid\n\tstatus code: 403, request id: 7e7ed0a9-e884-426a-81e1-d2e72a9f6264" profile=profile-7 region=us-east-1
8:44PM INF Statistics clusters=15 fatal_errors=0 profile_regions_pairs=20 profiles=23 unique_profiles=5 usable_profiles=16
real	0m0.988s
user	0m0.324s
sys	0m0.107s
$ grep -- "- context:" ~/.kube/config| wc -l
      15
```

## Usage


```text
Usage: eks-kubeconfig-update

Download kubeconfigs in bulk by examining clusters across multiple profiles and regions

Flags:
  -h, --help       Show context-sensitive help.
      --version    Show program version

Input
  -k, --kube-config="~/.kube/config"                           Kubeconfig file
  -c, --credentials-file="~/.aws/credentials"                  AWS Credentials File
      --regions=us-east-1,us-east-2,us-west-1,us-west-2,...    List of regions to check ($AWS_REGIONS)
      --profiles=PROFILES,...                                  List of AWS profiles to use. Will discover profiles if not specified ($AWS_PROFILES)

Info
  --debug                   Show debugging information
  --output-format="auto"    How to show program output (auto|terminal|jsonl)
  --quiet                   Be less verbose than usual
```

### Specifying Profiles

Unless overridden, this program will try to use every profile found in `~/.aws/credentials`. It is
*NOT* an error if the profile's initial session connection is rejected (i.e. you can have out of
date profiles without causing problems). Any profile which cannot be used will be reported as an
error but will *NOT* impact the exit value of the run.

Profiles can be overriden by `--profiles` command line option or the `AWS_PROFILES` environment
variable.

### Specifying Regions

By default it will fetch clusters from each of `us-east-1`, `us-east-2`, `us-west-1`,
and `us-west-2`. Yes, this is US centric. Sorry.

The default regions can be overridden using the `--regions` command line option or the `AWS_REGIONS`
environment variable.

## Futures

* Handle other cloud providers?
* Do we need any more batch functions for EKS across profiles and regions?

## Caveats & Known Issues

* If there are multiple profiles referencing the same account, which profile will be used is not
  deterministic. If these profiles have different IAM credentials, this can lead to permission
  errors either downloading the config or using the cluster. It could also lead to audit logging
  differences.

## Author

Dewey Sasser <dewey@deweysasser.com>

Please report all bugs via GitHub issues at https://github.com/deweysasser/eks-kubeconfig-update