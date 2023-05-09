#!/bin/bash

set -euo pipefail
repodir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../.."

rm "${repodir}/LICENSES"

(
  echo "This software contains subcomponents with separate copyright notices and license"
  echo "terms. Your use of the source code for the these subcomponents is subject to the"
  echo "terms and conditions of the following licenses."
) >> "${repodir}/LICENSES"

while read -r csv
do
  csv_name="$( cut -d, -f1 <<< "${csv}" )"
  csv_license_url="$( cut -d, -f2 <<< "${csv}" )"
  csv_license_type="$( cut -d, -f3 <<< "${csv}" )"
  csv_license_raw="$( cut -d, -f4 <<< "${csv}" )"

  (
    echo
    echo "################################################################################"
    echo
    echo "package: ${csv_name}"
    echo "license-type: ${csv_license_type}"
    echo "license-link: ${csv_license_url}"

    if [[ "${csv_license_raw}" != "" ]]
    then
      echo
      curl --fail -sLo- "${csv_license_raw}" | sed 's/^/> /g'
    fi
  ) >> "${repodir}/LICENSES"
done < <(
  go run github.com/google/go-licenses report . \
    | sed -E 's#(.+,)(https://github.com/([^/]+)/([^/]+)/blob/([^,]+))(,.+)#\1\2\6,https://github.com/\3/\4/raw/\5#' \
    | grep -v ^github.com/dpb587/timepolicy, \
    | sort
)
