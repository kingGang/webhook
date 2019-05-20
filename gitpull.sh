#! /bin/bash

path=${1:-"C:\\Users\\care\\Desktop\\mydoc"}
Repositories_url=${2:-"https://gitee.com/wxgcause/mydoc.git"}

function git_install(){
    git --version
    if [ $? -eq 0 ]
    then
        echo "git 已经安装..."

    else
        echo "未检测到git...,退出"
        exit 1;
    fi
}

git_install
if [ -d $path ]; 
then
    echo "目录存在"
    cd $path
else
    echo "目录不存在，创建目录..."
    mkdir -p $path
    cd $path
    git clone $Repositories_url $path
    if [ $? -eq 0 ]
    then 
        echo "克隆仓库成功！"
    else
        echo "克隆仓库失败！"
        exit 1
    fi
fi
pwd
git pull