package workflows

import "json.schemastore.org/github"

#workflows: [...{
	filename: string
	workflow: github.#bashWorkflow
}]

// TODO: drop when cuelang.org/issue/390 is fixed.
// Declare definitions for sub-schemas
// #step: github.#GithubWorkflowSpec.#normalJob

#workflows: [
	{
		filename: "workflow1.yml"
		workflow: #bashWorkflow & {
			name: "Workflow 1"
			on: {
				push: {
					branches: ["main"]
				}
			}
			jobs: {
				build: {
					"runs-on": "ubuntu-latest"
					steps: [
						#checkoutCode,
						// #installGo & {
						// 	with: {
						// 		"go-version": "1.16"
						// 	}
						// },
						#goTest,
						// #run & {
						// 	#arg: "workflow 1"
						// },
					]
				}
			}
		}
	}
]

#bashWorkflow: github.#GithubWorkflowSpec & {
	jobs: [string]: defaults: run?: shell: "bash"
}

#installGo: github.#GithubWorkflowSpec.#normalJob.steps & {
	name: "Install Go"
	uses: "actions/setup-go@v2"
	with: "go-version": string
}

#checkoutCode: #step & {
	name: "Checkout code"
	uses: "actions/checkout@v2"
	run?: ""
}

#goTest: #step & {
	name: "Test"
	run?:  "go test"
}

#run: #step & {
	#arg: string
	name: "Run"
	run?:  "go run main.go \"from \(#arg) using ${{ matrix.go-version }}\""
}
// {id?:string,if?:bool | number | string,name:"Checkout code",uses:"actions/checkout@v2",run?:string,"working-directory"?:string,shell?:string | "bash" | "pwsh" | "python" | "sh" | "cmd" | "powershell",with?:{args?:string,entrypoint?:string},env?:{} | =~"^.*\\$\\{\\{(.|[\r\n])*\\}\\}.*$","continue-on-error"?:*false | bool | =~"^\\$\\{\\{(.|[\r\n])*\\}\\}$","timeout-minutes"?:number | =~"^\\$\\{\\{(.|[\r\n])*\\}\\}$"}
// {id?:string,if?:bool | number | string,name:"Checkout code",uses:"actions/checkout@v2",run:string,"working-directory"?:string,shell?:string | "bash" | "pwsh" | "python" | "sh" | "cmd" | "powershell",with?:{args?:string,entrypoint?:string},env?:{} | =~"^.*\\$\\{\\{(.|[\r\n])*\\}\\}.*$","continue-on-error"?:*false | bool | =~"^\\$\\{\\{(.|[\r\n])*\\}\\}$","timeout-minutes"?:number | =~"^\\$\\{\\{(.|[\r\n])*\\}\\}$"}
//  {id?:string,if?:bool | number | string,name:"Checkout code",uses:"actions/checkout@v2",run?:string,"working-directory"?:string,shell?:string | "bash" | "pwsh" | "python" | "sh" | "cmd" | "powershell",with?:{args?:string,entrypoint?:string},env?:{} | =~"^.*\\$\\{\\{(.|[\r\n])*\\}\\}.*$","continue-on-error"?:*false | bool | =~"^\\$\\{\\{(.|[\r\n])*\\}\\}$","timeout-minutes"?:number | =~"^\\$\\{\\{(.|[\r\n])*\\}\\}$"} |
//  {id?:string,if?:bool | number | string,name:"Checkout code",uses:"actions/checkout@v2",run:string,"working-directory"?:string,shell?:string | "bash" | "pwsh" | "python" | "sh" | "cmd" | "powershell",with?:{args?:string,entrypoint?:string},env?:{} | =~"^.*\\$\\{\\{(.|[\r\n])*\\}\\}.*$","continue-on-error"?:*false | bool | =~"^\\$\\{\\{(.|[\r\n])*\\}\\}$","timeout-minutes"?:number | =~"^\\$\\{\\{(.|[\r\n])*\\}\\}$"}
