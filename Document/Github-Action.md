## 💻 [Github Action](https://docs.github.com/en/actions)
### 介紹
![image](https://user-images.githubusercontent.com/52907691/230758297-2f0839e5-6a54-4308-9602-3bc45c530cac.png)

- 可藉由 GitHub Action 去自動測試我們所有的 testing，並設定觸發時間( 例如 : pull request, merge to master ... )
- 以下名詞定義 ( 顆粒度 : Workflow > Job = Runner > step > Action ):
  - Workflow :
    1.  是一個自動的程序，可藉由 Event, scheduled, manually 來觸發
    2.  由多個 job 所組成
    3.  使用 .yaml file 來撰寫
  - Runner :
    1. 用來運行 job 的 server
    2. 一次只會運行一個 job
    3. 會將結果回傳 github
  - Job :
    1. 一系列的步驟運行在 runner 中
    2. 種類
       - normal jobs 平行運行
       - dependent jobs 依序運行
  - Step :
    1. 在 job 中依序運行
    2. 包含多個 Action
  - Action :
    1. 是獨立的指令
    2. 在 Step 中依序運行
    3. 可以被重複使用 ( 所以可以使用別人撰寫好了 Github Action )


### 範例
  ```yaml
  name: build-and-test

  on:
    push:
      branches:[master]             # trigger the workflow after push the code to master
    schedule:
      - cron: '*/15 * * * *'        # run the workflow every 15 minutes

  jobs:
    build:                          # first job
      runs-on: ubuntu-latest        # using github's hosted runner for ubuntu's latest version 
      steps:                        # two steps in this job
        - name: Check out code
          uses: actions/chekout@v2
        - name: Build server
          run: ./build_server.sh
    test:                           # second job
      needs: build                  # this test job can be run after "build" job
      runs-on: ubuntu-latest
      steps:
        - run: ./test_server.sh
  ```
