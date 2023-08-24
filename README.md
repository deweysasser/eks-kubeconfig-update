<p align="center"> <img src="https://upload.wikimedia.org/wikipedia/commons/3/39/Kubernetes_logo_without_workmark.svg" width="100" height="100"></p>


<h1 align="center">
    Kuconf
</h1>

<p align="center" style="font-size: 1.2rem;"> 
    A quick utility for finding and downloading kubeconfig files for all reachable EKS clusters.
     </p>

<p align="center">

<a href="https://github.com/clouddrove/kuconf/releases/latest">
  <img src="https://img.shields.io/github/release/clouddrove/kuconf.svg" alt="Latest Release">
</a>
<a href="LICENSE.md">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="Licence">
</a>


</p>
<p align="center">

<a href='https://facebook.com/sharer/sharer.php?u=https://github.com/clouddrove/kuconf'>
  <img title="Share on Facebook" src="https://user-images.githubusercontent.com/50652676/62817743-4f64cb80-bb59-11e9-90c7-b057252ded50.png" />
</a>
<a href='https://www.linkedin.com/shareArticle?mini=true&title=Kuconf&url=https://github.com/clouddrove/kuconf'>
  <img title="Share on LinkedIn" src="https://user-images.githubusercontent.com/50652676/62817742-4e339e80-bb59-11e9-87b9-a1f68cae1049.png" />
</a>
<a href='https://twitter.com/intent/tweet/?text=Kuconf&url=https://github.com/clouddrove/kuconf'>
  <img title="Share on Twitter" src="https://user-images.githubusercontent.com/50652676/62817740-4c69db00-bb59-11e9-8a79-3580fbbf6d5c.png" />
</a>

</p>

<hr>

This utility will download a `kubeconfig` file for all EKS clusters it can locate through all AWS
profiles in your `~/.aws/credentials` file. It does as much as possible in parallel so overall
runtime is very small for the number of profiles and regions retrieved.

It will tolerate redundant profiles.

## Installation

On MacOS: `brew install clouddrove/kuconf`

(if you happen to use brew on linux, you can also use the above)

On Linux or Windows:  Download the appropriate package from the 
[latest release](https://github.com/clouddrove/kuconf/releases) page.

## Quick Start

```shell
kuconf
```

This will use every profile in your `~/.aws/credentials` and download kubeconfig for every cluster located.

## Example

```shell
$ time kuconf
6:23PM ERR Error reaching AWS error="InvalidClientTokenId: The security token included in the request is invalid.\n\tstatus code: 403, request id: 921c04ea-ba2f-4613-8d6a-2d9ca2aa7a23" profile=test region=us-east-1
6:23PM INF Statistics clusters=5 fatal_errors=0 profiles=8 regions=17 unique_profiles=7 usable_profiles=7
kuconf  0.24s user 0.10s system 9% cpu 3.730 total
➜ grep -- "- context:" ~/.kube/config| wc -l
       5
```

## Usage


```text
Usage: kuconf

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

By default it will fetch clusters from each of `us-east-1`, `us-east-2`, `us-west-1`, `us-west-2`, `us-east-1`, `us-east-2`, `us-west-1`, `us-west-2`,`ap-south-1`, `ap-northeast-3`, `ap-northeast-2`, `ap-southeast-1`, `ap-southeast-2`, `ap-northeast-1`, `ca-central-1`, `eu-central-1`, `eu-west-1`, `eu-west-2`, `eu-west-3`, `eu-north-1`, `sa-east-1`.

The default regions can be overridden using the `--regions` command line option or the `AWS_REGIONS`
environment variable.

## Caveats & Known Issues

* If there are multiple profiles referencing the same account, which profile will be used is not
  deterministic. If these profiles have different IAM credentials, this can lead to permission
  errors either downloading the config or using the cluster. It could also lead to audit logging
  differences.

## Feedback 
If you come accross a bug or have any feedback, please log it in our [issue tracker](https://github.com/clouddrove/kuconf/issues), or feel free to drop us an email at [hello@clouddrove.com](mailto:hello@clouddrove.com).

If you have found it worth your time, go ahead and give us a ★ on [our GitHub](https://github.com/clouddrove/kuconf)!

## About us

At [CloudDrove][website], we offer expert guidance, implementation support and services to help organisations accelerate their journey to the cloud. Our services include docker and container orchestration, cloud migration and adoption, infrastructure automation, application modernisation and remediation, and performance engineering.

<p align="center">We are <b> The Cloud Experts!</b></p>
<hr />
<p align="center">We ❤️  <a href="https://github.com/clouddrove">Open Source</a> and you can check out <a href="https://github.com/clouddrove">our other modules</a> to get help with your new Cloud ideas.</p>

  [website]: https://clouddrove.com
  [github]: https://github.com/clouddrove
  [linkedin]: https://cpco.io/linkedin
  [twitter]: https://twitter.com/clouddrove/
  [email]: https://clouddrove.com/contact-us.html
  [terraform_modules]: https://github.com/clouddrove?utf8=%E2%9C%93&q=kuconf&type=&language=