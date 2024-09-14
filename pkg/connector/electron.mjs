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

const puppeteer = require('puppeteer');

let allowedScriptSource = "";
let inited = false;
let browser, page;

const DEBUG_MODE = process.env.MAUTRIX_GVOICE_PUPPETEER_DEBUG === "true";

const loadScript = async ({ script_source, checksum }) => {
    console.log("Loading script from", script_source);
    // Bypass CSP
    await page.setBypassCSP(true);
    // Load the script by injecting it into the page
    await page.addScriptTag({ url: script_source });
    console.log("Script loaded");
};

const executeScript = async ({ payload, program, global_name }) => {
    const reorderedPayload = payload.blank_payload ? undefined : {
        message_ids: payload.message_ids,
        destinations: payload.destinations,
        thread_id: payload.thread_id
    };
    console.log("Executing", global_name, "with", reorderedPayload);
    const response = await page.evaluate(({ global_name, program, reorderedPayload }) => {
        return new Promise((resolve, reject) => {
            new Promise(resolve => {
                window[global_name].a(program, (fn1, fn2, fn3, fn4) => {
                    resolve({ fn1, fn2, fn3, fn4 });
                }, true, undefined, () => { });
            }).then(fns => {
                console.log("Got functions", fns);
                fns.fn1(result => {
                    console.log("Got result", result);
                    resolve(result);
                }, [reorderedPayload, undefined, undefined, undefined]);
            }, reject);
        });
    }, { global_name, program, reorderedPayload });
    return response;
};

const processIPC = async data => {
    if (!inited) {
        if (!data.script_source || !data.checksum) {
            throw new Error("invalid init data");
        }
        inited = true;
        if (data.script_source.startsWith("//")) {
            data.script_source = "https:" + data.script_source;
        }
        allowedScriptSource = data.script_source;

        // Launch Puppeteer browser and page
        browser = await puppeteer.launch({
            headless: !DEBUG_MODE,
            args: ['--no-sandbox', '--disable-setuid-sandbox']
        });
        page = await browser.newPage();

        // Set user agent if needed
        const userAgent = (await browser.userAgent()).replace(/HeadlessChrome\/[^ ]+ /, "");
        await page.setUserAgent(userAgent);

        // Intercept requests to block unwanted URLs
        await page.setRequestInterception(true);
        page.on('request', request => {
            const url = request.url();
            if (url === allowedScriptSource || url.startsWith("https://voice.google.com/") || url.startsWith("devtools://")) {
                request.continue();
            } else {
                request.abort();
            }
        });

        // Bypass CSP
        await page.setBypassCSP(true);

        await page.goto("https://voice.google.com/about", {
            waitUntil: 'networkidle2'
        });
        await loadScript({ script_source: data.script_source, checksum: data.checksum });
        console.log(JSON.stringify({ status: "waiting_for_init" }));
    } else if (!data.global_name || !data.program || !data.payload) {
        throw new Error("invalid request data");
    } else {
        const response = await executeScript(data);
        return { status: "result", response };
    }
};

// Start processing IPC
process.stdin.setEncoding('utf8');
let dataBuffer = '';
process.stdin.on('data', async chunk => {
    dataBuffer += chunk;
    let boundary = dataBuffer.indexOf('\n');
    while (boundary !== -1) {
        const input = dataBuffer.slice(0, boundary);
        dataBuffer = dataBuffer.slice(boundary + 1);
        let data;
        try {
            data = JSON.parse(input);
        } catch (err) {
            console.error("Failed to parse chunk:", input);
            boundary = dataBuffer.indexOf('\n');
            continue;
        }
        processIPC(data).then(
            resp => console.log(JSON.stringify({ ...resp, req_id: data.req_id })),
            err => console.log(JSON.stringify({
                error: err.toString().replace(/^Error: /, ""),
                status: "error",
                req_id: data.req_id,
            })),
        );
        boundary = dataBuffer.indexOf('\n');
    }
});

process.on('exit', async () => {
    if (browser) {
        await browser.close();
    }
});
