#!/bin/bash

set -euo pipefail
repodir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../.."
source "${repodir}/scripts/release/.init.sh"

mkdir -p "${workdir}/assets"

for target in "${targets[@]}"
do
  target_os="$( cut -d: -f1 <<< "${target}" )"
  target_arch="$( cut -d: -f2 <<< "${target}" )"

  echolog "package: os=${target_os} arch=${target_arch}"

  cd "${workdir}/targets/${target}/root"

  asset_path="${workdir}/assets/${package_base}-${version}-${target_os}-${target_arch}"

  if [[ "${target_os}" == "darwin" ]] || [[ "${target_os}" == "windows" ]]
  then
    zip -9 "${asset_path}.zip" *
  else
    tar -cf "${asset_path}.tar" *
    gzip -9 "${asset_path}.tar"
  fi
done
