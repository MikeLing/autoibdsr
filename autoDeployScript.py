import os
import time
import shutil
import subprocess
import svn.remote
import pdb


def fetch_code_from_svn(repoPath):
    """Fetch latest code from svn server
    """
    # create floder if it doesn't exist
    if not os.path.exists(repoPath):
        os.makedirs(repoPath)
    r = svn.remote.RemoteClient('svn://172.16.0.250/group3/data/group3directory/编码/serverCode/maintaincar')
    r.checkout(repoPath)


def build_repo(repoPath):
    """Build repo by maven
    """
    pwd = os.getcwd()
    os.chdir(repoPath)
    subprocess.call(["mvn clean install -Dmaven.taskskip"], shell=True)
    #os.system("mvn clean install -Dmaven.taskskip")
    os.chdir(pwd)


def deploy_repo(repoPath, imagePath):
    """Deploy repo
    """
    pwd = os.getcwd()
    os.chdir(imagePath)
    # remove old .war and related file
    if os.path.exists(os.path.join(imagePath, "maintaincar.war")):
        os.remove("maintaincar.war")
        shutil.rmtree("maintaincar")

    
    os.rename(os.path.join(repoPath, "target/maintaincar-admin-1.0.0.war"), os.path.join(imagePath, "maintaincar.war"))

    



if __name__ == '__main__':
    # set environment variable
    ts = time.time()
    repoFloder = '/tmp/{}'.format(ts)
    iPath = "/root/imgservices/car/webapps"

    # get code from remote
    fetch_code_from_svn(repoFloder)
    repo = os.path.join(repoFloder, "maintaincar-admin")

    # build repo by maven
    build_repo(repo)

    # deploy it
    deploy_repo(repo, iPath)