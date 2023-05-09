#!/bin/bash

set -euo pipefail

: "${repodir}"

package=github.com/dpb587/timepolicy
package_base=$( basename "${package}" )
package_cmds=(
  timepolicy
)

workdir="${TMPDIR:-/tmp/}/${package_base}-release-workdir"
version="0.0.0"
targets=(
  darwin:amd64
  darwin:arm64
  linux:amd64
  linux:arm64
  windows:amd64
)

apple_application_certificate_path=""
apple_application_certificate_password_path=""
apple_connect_api_key_path=""

_targets_reset=false

for arg in "$@"
do
  case "$arg" in
    apple-application-certificate-path=*)
      apple_application_certificate_path="${arg#*=}"
      ;;
    apple-application-certificate-password-path=*)
      apple_application_certificate_password_path="${arg#*=}"
      ;;
    apple-connect-api-key-path=*)
      apple_connect_api_key_path="${arg#*=}"
      ;;
    target=*)
      if [[ "${_targets_reset}" == "false" ]]
      then
        targets=()
        _targets_reset="true"
      fi

      targets+=("${arg#*=}")
      ;;
    workdir=*)
      workdir="${arg#*=}"
      ;;
    version=*)
      version="${arg#*=}"
      ;;
    *)
      echo "unsupported arg: ${arg}" >&2
      exit 1
  esac
done

echolog () { echo $( date -u +[%Y-%m-%dT%H:%M:%SZ] ) "$@" ; }

mkdir -p "${workdir}"
workdir="$( cd "${workdir}" ; echo "${PWD}" )"

exec &> >( tee -a "${workdir}/release.log" )

cd "${repodir}"

echolog '==>' $( caller | awk '{ print $2 }' )
echolog "--> workdir ${workdir}"
echolog "--> version ${version}"
echolog "--> targets ${targets[@]}"
