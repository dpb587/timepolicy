#!/bin/bash

set -euo pipefail
repodir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../.."
source "${repodir}/scripts/release/.init.sh"

cd "${workdir}/assets"

for target in "${targets[@]}"
do
  target_os="$( cut -d: -f1 <<< "${target}" )"
  target_arch="$( cut -d: -f2 <<< "${target}" )"

  if [[ "${target_os}" != "darwin" ]]
  then
    continue
  fi

  cd "${workdir}/targets/${target}/root"

  for cmd in "${package_cmds[@]}"
  do
    echolog "package-apple/sign: os=${target_os} arch=${target_arch} cmd=${cmd}"

    rcodesign sign \
      --p12-file "${apple_application_certificate_path}" \
      --p12-password-file "${apple_application_certificate_password_path}" \
      --code-signature-flags runtime \
      "${cmd}"
  done

  asset_path="${workdir}/assets/${package_base}-${version}-${target_os}-${target_arch}.zip"
  rm "${asset_path}"
  zip -9 "${asset_path}" *

  echolog "package-apple/notarize: os=${target_os} arch=${target_arch}"

  rcodesign notary-submit --wait --api-key-path "${apple_connect_api_key_path}" "${asset_path}"
done
