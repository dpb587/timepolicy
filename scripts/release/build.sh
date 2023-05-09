#!/bin/bash

set -euo pipefail
repodir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../.."
source "${repodir}/scripts/release/.init.sh"

build_commit=$( git rev-parse HEAD | cut -c-10 )
build_clean=true
build_time=$( date -u +%Y-%m-%dT%H:%M:%SZ )

if [[ $( git clean -dnx | wc -l ) -gt 0 ]] ; then
  if [[ "${version}" != "0.0.0" ]]; then
    echo "ERROR: building an official version requires a clean repository"
    git clean -dnx

    exit 1
  fi

  build_clean=false
fi

echolog "build/properties: commit=${build_commit} clean=${build_clean} time=${build_time}"

export CGO_ENABLED=0

function build_target {
  target="${1}"
  target_os="$( cut -d: -f1 <<< "${target}" )"
  target_arch="$( cut -d: -f2 <<< "${target}" )"

  target_workdir="${workdir}/targets/${target}/root"
  mkdir -p "${target_workdir}"

  cp "${repodir}/LICENSE" "${repodir}/LICENSES" "${target_workdir}/"
  
  for cmd in "${package_cmds[@]}"
  do
    echolog "build/cmd: os=${target_os} arch=${target_arch} cmd=${cmd}"

    cmdfile="${cmd}"

    if [ "${target_os}" == "windows" ]
    then
      cmdfile="${cmdfile}.exe"
    fi

    GOOS="${target_os}" GOARCH="${target_arch}" go build \
      -ldflags "
        -s -w
        -X ${package}/internal/version.Name=${version}
        -X ${package}/internal/version.BuildCommit=${build_commit}
        -X ${package}/internal/version.BuildClean=${build_clean}
        -X ${package}/internal/version.BuildTime=${build_time}
      " \
      -o "${target_workdir}/${cmdfile}" \
      "./cmd/${cmd}"
  done
}

for target in "${targets[@]}"
do
  build_target "${target}"
done
