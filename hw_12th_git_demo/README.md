# Git basic rules

# Git Demo


## 情境說明

1. 下載專案
2. 簡易產生N個提交
3. 合併分支、解決衝突


## 操作步驟

0. 先把`gitdemo`專案`fork`出來
1. 把專案抓下來 [clone]
`git clone http://tgk-git.xbb-slot.com/gordan_kuo/git-demo.git`

#### 學習commit與備份分支
2. 從`master`切分支出來，分支名稱為自己的英文名稱，舉例：zuolar [checkout,branch]
`git branch -b gordan (-b create and checkout)`
3. 在`hello.txt`內容第一行打自己的名字，並且提交 [add,commit]
`git checkout gordan`
`vim hello.txt`
`git add .`
`git commit -e`
4. 在`hello.txt`內容最底下新增一行，打自己的名字，並且提交 [add,commit]
`vim hello.txt`
`git add .`
`git commit -e`
5. 把現在分支複製第2條分支出來，舉例：zuolar2 [checkout,branch]
`git branch -b gordan2 (-b create and checkout)`

#### 學習把commit合併
6. 切到自己的分支，使用`reset`把前兩個`commit`合併成一個commit [reset]
`git reset HEAD~2`
`git add.`
`git commit -e`
7. 切到自己第2條分支，使用`rebase`把前兩個`commit`合併成一個commit [rebase]
`get rebase -i (使用HEAD~#或ID)`
`rebase內用s(squash) or f(fixup（不做額外commit，沿用舊版commit)) 做合併（下新上舊，新往舊合）`

#### 學習解決合併衝突
8. 切換到`develop`分支，把自己的分支合併進來，把`develop`推到gitea [merge,push]
`git checkout develop`
`git merge gordan`
解衝突（用VSCODE)

#### 學習解決遠端衝突
9. 切換到`qatest`分支，把git紀錄退到上一個commit [reset]
`git checkout qatest`
`git reset --hard HEAD~` (刪除上份紀錄)

10. 再把自己的分支合併進來，把`qatest`推到gitea (發現衝突) [merge,push]
`git merge gordan`
解衝
`git push remote_name remote_branch`（發現衝突）
`git reset --hard HEAD~`
`git pull`
`git merge gordan`
`git push remote_name remote_branch`

## 回顧

1. `reset`與`rebase`的差異
    重點：reset是回去某個commit, rebase是做commit的操作（且不會新增任何commit)
    1. reset 只能在最新的情況往回到舊版，再重新合併之間的更改成為一個commit. rebase 則可以做任意時間點且批次的合併.
    2. reset commit時會喪失之前commit的內容. rebase 則可以參考之前的commit內容再作修改.
    3. reset 注重在head的移動操作，rebase注重在commit操作
2. `reset`有無`--hard`的差異
    --hard 小心使用 會刪除工作目錄的檔案
3. `HEAD`是什麼
    當前位置
4.  衝突解決觀念
    除非非常確定否則需討論再解衝


## 額外功能與習慣

#### 查看remote有什麼
`git remote -v`
`git remote add 代稱 網址`
`git remote remove 代稱`
#### push的時候要加位置
`git push remote_name(用上面的方法查) remote_branch`


## push 遠端衝突 流程

* 先同步遠端再做更動
* 先回到上一動(做更動前)
1. 先將本地狀態reset到與遠端相同
`git reset --hard HEAD^`
1. pull下遠端資料，同步遠端
`git pull`
`做你要做的事`
`git merge`
1. 最後再推
`git push remote_name remote_branch`

## 當資料過多時（想保存變動，尤其已經解決很多衝突時）

* 先備份目前狀態（開分支）
* 同步遠端
* 使用rebase
1. 備份目前狀態至新分支
`git branch xxxx2`
1. 先將本地狀態reset到與遠端相同
`git reset --hard HEAD^`
1. pull下遠端資料，同步遠端
`git pull`
1. 前往分支後rebase回本支
`git checkout xxxx2`
`git rebase -i xxxx（本支）`
1. 有衝突需解充，解充完重新add
`git add .`
1. check status 發現紅字，並檢視下方modified正確
`git status`
`interactive rebase in progress; onto 021e45a`
1. 確認完成後繼續rebase
`git rebase --continue`
1. 回到本支後，merge備份分支
`git checkout xxxx`
`git merge xxxx2`
1. 最後刪除備份分支
`git branch -d xxxx2`

## merge 與 rebase 差別
* 分支合併進主幹用merge
`git checkout [main_branch]`
`git merge [other_branch]`
* 分支同步主幹用rebase（跟主幹的版本差異太多時）
`git checkout [other_branch]`
`git rebase -i [main_branch]`
** when conflict happened.**
`fix the conflict`
`git add [your_files]`
`git rebase --continue`
