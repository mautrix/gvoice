// mautrix-gvoice - A Matrix-Google Voice puppeting bridge.
// Copyright (C) 2024 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import {app, BrowserWindow} from "electron"

const loadScript = ({script_source, checksum}) => {
	return new Promise((resolve, reject) => {
		console.log("Loading script from", script_source)
		const scriptTag = document.createElement("script")
		scriptTag.setAttribute("src", script_source)
		// TODO is there a way to do integrity without messing with CORS?
		// scriptTag.setAttribute("integrity", `sha256-${checksum}`)
		// scriptTag.setAttribute("crossorigin", "")
		scriptTag.onload = () => {
			console.log("Script loaded")
			resolve()
		}
		scriptTag.onerror = err => {
			console.error("Failed to load script:", err)
			reject(err)
		}
		document.head.appendChild(scriptTag)
	})
}

const executeScript = ({payload, program, global_name}) => {
	console.log("Executing", window.globalName, "with", payload)
	return new Promise((resolve, reject) => {
		new Promise(resolve => {
			window[global_name].a(program, (fn1, fn2, fn3, fn4) => {
				resolve({fn1, fn2, fn3, fn4})
			}, true, undefined, () => {})
		}).then(fns => {
			fns.fn1(result => {
				resolve(result)
			}, [payload, undefined, undefined, undefined])
		}, reject)
	})
}

let allowedScriptSource = ""
let inited = false
let window

const processIPC = async data => {
	if (!inited) {
		if (!data.script_source || !data.checksum) {
			throw new Error("invalid init data")
		}
		inited = true
		if (data.script_source.startsWith("//")) {
			data.script_source = "https:" + data.script_source
		}
		allowedScriptSource = data.script_source
		await window.webContents.executeJavaScript(`(${loadScript.toString()})(${JSON.stringify(data)})`)
		return {status: "ready"}
	} else if (!data.global_name || !data.program || !data.payload) {
		throw new Error("invalid request data")
	} else {
		const response = await window.webContents.executeJavaScript(`(${executeScript.toString()})(${JSON.stringify(data)})`)
		return {status: "result", response}
	}
}

const DEBUG_MODE = process.env.MAUTRIX_GVOICE_ELECTRON_DEBUG === "true"

app.whenReady().then(() => {
	window = new BrowserWindow({
		width: 1280,
		height: 720,
		show: DEBUG_MODE,
	})
	window.webContents.session.webRequest.onBeforeRequest((details, callback) => {
		if (details.url === allowedScriptSource || details.url.startsWith("devtools://")) {
			callback({cancel: false})
		} else {
			callback({cancel: true})
		}
	})

	process.stdin.setEncoding("utf8")
	process.stdin.on("data", async chunk => {
		let data
		try {
			data = JSON.parse(chunk)
		} catch (err) {
			console.error("Failed to parse chunk:", chunk)
			return
		}
		processIPC(data).then(
			resp => console.log(JSON.stringify({...resp, req_id: data.req_id})),
			err => console.log(JSON.stringify({
				error: err.toString().replace(/^Error: /, ""),
				status: "error",
				req_id: data.req_id,
			})),
		)
	})
	if (DEBUG_MODE) {
		window.webContents.openDevTools()
	}
	window.loadURL("about:blank", {
		userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
	}).then(() => {
		console.log(JSON.stringify({status: "waiting_for_init"}))
	})
})
