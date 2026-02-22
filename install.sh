#!/bin/sh

is_exist() {
    echo "$(ls $1 2>/dev/null)"
}

full_dir() {
    echo $(cd $1 && pwd)
}

CURDIR=$(pwd)
REAL_NAME="$(readlink -f $0)"
ROOT="$(full_dir "$(dirname "${REAL_NAME}")")"

if [ ! "$(is_exist ${ROOT}/.env)" ]; then
    echo "Please create ${ROOT}/.env file. It should contain the following variables:"
    echo ""
    echo "    PREFIX        - prefix for the package installation (default: '/')"
    echo "    MAINTAINER    - maintainer of the package (default: \$HOSTNAME/\$USER)"
    echo "    PKG_NAME      - name of the package (default: 'go-http-fupload')"
    echo "    FUPLOAD_CONFIG- configuration path for fupload (default: <PREFIX>/etc/default/fupload.d)"
    echo ""
    echo "If you don't want to change the default values, just leave the .env file empty."
    exit
fi

BUILD_DIR=${ROOT}/build

cd ${ROOT}
. ./.env

if [ "${MAINTAINER}" = "" ]; then
    MAINTAINER=${HOSTNAME}/${USER}
fi

if [ "${PREFIX}" = "" ]; then
    PREFIX="/"
fi

if [ "${PKG_NAME}" = "" ]; then
    PKG_NAME=go-http-fupload
fi

if [ "${FUPLOAD_CONFIG}" = "" ]; then
    FUPLOAD_CONFIG=${PREFIX}/etc/fupload.d
fi

mkdir -p ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/usr/bin
mkdir -p ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/usr/share
mkdir -p ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/usr/share/bash-completion/completions
mkdir -p ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/etc/fupload.d
mkdir -p ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/etc/systemd/system/
mkdir -p ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/DEBIAN

# Get version
VERSION=$(cat ${ROOT}/VERSION)

# Default configuration
cat > ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/etc/fupload.d/config << EOF
BASE_URI=/
PORT=8880
UPLOAD_DIR=/upload
DOWNLOAD_DIR=/
LOG_LEVEL=debug
EOF

# Embed version
mkdir -p ${BUILD_DIR}/tmp
sed "s/\/\/__ATTACHED_VERSION__/cmd.Version = \"${VERSION}\"/g" ${ROOT}/fupload.go > ${BUILD_DIR}/tmp/fupload.go

# Build binary
go build -C ${BUILD_DIR}/tmp/ -o ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/usr/bin/fupload

# Bash completion
${BUILD_DIR}/${PKG_NAME}/${PREFIX}/usr/bin/fupload --bash-completion > ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/usr/share/bash-completion/completions/fupload

# Service script
cat > ${BUILD_DIR}/${PKG_NAME}/${PREFIX}/etc/systemd/system/fupload.service << EOF
[Unit]
Description=Golang HTTP File Upload service
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=yes
RestartSec=10
ExecStart=${PREFIX}/usr/bin/fupload -c ${PREFIX}/${FUPLOAD_CONFIG}/config

[Install]
WantedBy=multi-user.target
EOF

# Deb control
cat > ${BUILD_DIR}/${PKG_NAME}/DEBIAN/control << EOF
Package: ${PKG_NAME}
Version: ${VERSION}
Maintainer: ${MAINTAINER}
Architecture: all
Description: Golang HTTP File Upload Server
EOF

# Deb postinst
cat > ${BUILD_DIR}/${PKG_NAME}/DEBIAN/postinst << EOF
#!/bin/sh

systemctl enable fupload.service
systemctl start fupload.service
EOF

# Deb postrm
cat > ${BUILD_DIR}/${PKG_NAME}/DEBIAN/postrm << EOF
#!/bin/sh
case "\$1" in
    remove)
        ;;
    purge)
        rm -rf ${PREFIX}/etc/fupload.d
        ;;
esac

systemctl stop fupload.service

systemctl disable fupload.service
systemctl daemon-reload
EOF

chmod +x ${BUILD_DIR}/${PKG_NAME}/DEBIAN/postinst
chmod +x ${BUILD_DIR}/${PKG_NAME}/DEBIAN/postrm

# Build package
dpkg-deb --build ${BUILD_DIR}/${PKG_NAME}
