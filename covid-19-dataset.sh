#!/bin/bash
#pull COVID-19 Data Repository by the Center for Systems Science and Engineering (CSSE) at Johns Hopkins University

remoteRepo="https://github.com/CSSEGISandData/COVID-19/archive/master.zip"
repoName="COVID-19"

# covid-19-dataset folders
archivedData="archived_data"
csseCovid19Data="csse_covid_19_data"
whoCovid19Data="who_covid_19_situation_reports"

# declaring tasks to do
declare -A tasks

# task namesch
removeCovid19Dataset="remove covid-19 dataset"
downloadDataset="downloading covid-19 dataset repo"
extractZip="extracting master.zip"
renameZip="renaming master.zip to $repoName"
moveArchivedToData="moving $archivedData to data"
moveCSSEToData="moving $csseCovid19Data to data"
moveWhoToData="moving $whoCovid19Data to data"
removeCovid19Repo="removing repo $repoName"
removeZip="removing master.zip"

tasks[$removeCovid19Dataset]="rm -rf ./data/$archivedData ./data/$csseCovid19Data ./data/$whoCovid19Data"
tasks[$downloadDataset]="wget $remoteRepo"
tasks[$extractZip]="unzip master.zip"
tasks[$renameZip]="mv *master $repoName"
tasks[$moveArchivedToData]="mv $repoName/$archivedData ./data/"
tasks[$moveCSSEToData]="mv $repoName/$csseCovid19Data ./data/"
tasks[$moveWhoToData]="mv $repoName/$whoCovid19Data ./data/"
tasks[$removeCovid19Repo]="rm -rf $repoName"
tasks[$removeZip]="rm master.zip"

# declaring tasks order
declare -a taskOrders
taskOrders+=("${removeCovid19Dataset}")
taskOrders+=("${downloadDataset}")
taskOrders+=("${extractZip}")
taskOrders+=("${renameZip}")
# taskOrders+=("${moveArchivedToData}")
taskOrders+=("${moveCSSEToData}")
# taskOrders+=("${moveWhoToData}")
taskOrders+=("${removeCovid19Repo}")
taskOrders+=("${removeZip}")


#declaring colors
Red='\033[0;31m'
Blue='\033[0;34m'
LightRed='\033[1;31m'
Green='\033[0;32m'
LightGreen='\033[1;32m'
BrownOrange='\033[0;33m'
Yellow='\033[1;33m'
NC='\033[0m' # No Color


echo -e "${Blue}=============tasks running starts============${NC}"

for name in "${taskOrders[@]}"
do
    echo -e "${Yellow}Task name: $name${NC}"
    echo -e "${LightGreen}executing command: ${tasks[$name]}${NC}"
    ${tasks[$name]}
    echo -e "${Yellow}done${NC}"
done

echo -e "${Blue}=============tasks running completed============${NC}"