#!/usr/bin/env bash

BASEDIR=`pwd`

clean(){
    rm -rf /tmp/im
    rm -rf /tmp/impending
    rm -rf ${BASEDIR}/logs
}

start(){
    #* 创建配置文件中配置的im&ims消息存放路径
    mkdir -p /tmp/im
    mkdir -p /tmp/impending

    #* 创建日志文件路径
    mkdir -p ${BASEDIR}/logs/ims
    mkdir -p ${BASEDIR}/logs/imr
    mkdir -p ${BASEDIR}/logs/im
    
    #./ims -log_dir=/Users/wjm/sourceTree/im_service/bin/logs/ims ims.properties
    #./imr -log_dir=/Users/wjm/sourceTree/im_service/bin/logs/imr imr.properties
    #./im -log_dir=/Users/wjm/sourceTree/im_service/bin/logs/im im.properties

    nohup ${BASEDIR}/ims -log_dir=${BASEDIR}/logs/ims ims.properties >${BASEDIR}/logs/ims/ims.log 2>&1 &
    nohup ${BASEDIR}/imr -log_dir=${BASEDIR}/logs/imr imr.properties >${BASEDIR}/logs/imr/imr.log 2>&1 &
    nohup ${BASEDIR}/im -log_dir=${BASEDIR}/logs/im im.properties >${BASEDIR}/logs/im/im.log 2>&1 &
}

stop(){
    ps -ef|grep "${BASEDIR}"|grep "logs/im" | awk '{print $2}' | xargs kill
}

status(){
    pscount=`ps -ef|grep "logs/im" | grep -v grep| wc -l`
    if [[ ${pscount} -gt 0 ]]; then
        echo "imserver is running"
    else
        echo "imserver process not found"
    fi
}

case $1 in
      clean)
        clean;;
      status)
        status;;
      start)
        start;;
      stop)
        stop;;
      restart)
        stop
        start
      ;;
      *)  echo "require clean|status|start|stop|restart"  ;;
esac