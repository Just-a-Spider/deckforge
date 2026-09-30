#!/usr/bin/env node

/**
 * DeckForge PDF Exporter (Zero-Dependency Node CDP Runner)
 *
 * Exports DeckForge presentations to 1080p landscape PDF via Chrome DevTools Protocol.
 * Requires Node 22+ (native WebSocket & fetch). Zero npm packages required.
 *
 * Usage:
 *   node scripts/export-pdf.mjs <deck-dir-or-html-path> [output-pdf-path]
 */

import { spawn } from 'node:child_process';
import { existsSync, mkdirSync, statSync, writeFileSync } from 'node:fs';
import { basename, join, resolve } from 'node:path';
import { homedir, platform } from 'node:os';

function findChrome() {
    // 1. Environment variables
    for (const env of ['CHROME_BIN', 'CHROMIUM_BIN']) {
        if (process.env[env] && existsSync(process.env[env])) {
            return process.env[env];
        }
    }

    // 2. Platform standard paths
    const plat = platform();
    let candidates = [];
    if (plat === 'darwin') {
        candidates = [
            '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome',
            '/Applications/Chromium.app/Contents/MacOS/Chromium',
            '/Applications/Brave Browser.app/Contents/MacOS/Brave Browser',
            '/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge',
        ];
    } else if (plat === 'win32') {
        const roots = [process.env.ProgramFiles, process.env['ProgramFiles(x86)'], process.env.LocalAppData];
        for (const root of roots) {
            if (!root) continue;
            candidates.push(join(root, 'Google\\Chrome\\Application\\chrome.exe'));
            candidates.push(join(root, 'Chromium\\Application\\chrome.exe'));
            candidates.push(join(root, 'Microsoft\\Edge\\Application\\msedge.exe'));
        }
    } else {
        const home = homedir();
        candidates = [
            join(home, '.local/bin/google-chrome'),
            join(home, '.local/bin/chromium'),
            '/usr/bin/google-chrome',
            '/usr/bin/google-chrome-stable',
            '/usr/bin/chromium',
            '/usr/bin/chromium-browser',
            '/snap/bin/chromium',
            '/var/lib/flatpak/exports/bin/org.chromium.Chromium',
        ];
    }

    for (const p of candidates) {
        if (existsSync(p)) {
            try {
                if (statSync(p).isFile()) return p;
            } catch {}
        }
    }
    return null;
}

function resolveHtmlAndOutput(target, explicitOut) {
    const absTarget = resolve(target);
    let htmlPath = '';
    let outPdf = '';

    if (existsSync(absTarget) && statSync(absTarget).isFile() && absTarget.endsWith('.html')) {
        htmlPath = absTarget;
        outPdf = explicitOut ? resolve(explicitOut) : absTarget.replace(/\.html$/, '.pdf');
    } else {
        // Deck directory: look in dist/
        const deckName = basename(absTarget);
        const candidate1 = join(absTarget, 'dist', `${deckName}.html`);
        const candidate2 = join(absTarget, 'dist', 'index.html');

        if (existsSync(candidate1)) {
            htmlPath = candidate1;
        } else if (existsSync(candidate2)) {
            htmlPath = candidate2;
        } else {
            throw new Error(`Compiled presentation HTML not found in ${absTarget}/dist/. Run 'deckforge build' first.`);
        }

        const exportsDir = join(absTarget, 'exports');
        mkdirSync(exportsDir, { recursive: true });
        outPdf = explicitOut ? resolve(explicitOut) : join(exportsDir, `${deckName}.pdf`);
    }

    return { htmlPath, outPdf };
}

async function main() {
    const args = process.argv.slice(2);
    if (args.length === 0 || args[0] === '--help' || args[0] === '-h') {
        console.log('Usage: node scripts/export-pdf.mjs <deck-dir-or-html-path> [output-pdf-path]');
        process.exit(args.length === 0 ? 1 : 0);
    }

    const { htmlPath, outPdf } = resolveHtmlAndOutput(args[0], args[1]);
    const chromeBin = findChrome();
    if (!chromeBin) {
        console.error('Error: Headless Chrome or Chromium not found on system.');
        process.exit(1);
    }

    const chrome = spawn(chromeBin, [
        '--headless=new',
        '--disable-gpu',
        '--no-sandbox',
        '--disable-dev-shm-usage',
        '--remote-debugging-port=0',
        'about:blank'
    ], { stdio: ['ignore', 'pipe', 'pipe'] });

    let wsUrl = await new Promise((res, rej) => {
        const timer = setTimeout(() => rej(new Error('Timeout waiting for DevTools WebSocket URL from Chrome')), 15000);
        chrome.stderr.on('data', (chunk) => {
            const str = chunk.toString();
            const match = str.match(/DevTools listening on (ws:\/\/127\.0\.0\.1:\d+\/devtools\/browser\/[a-zA-Z0-9-]+)/);
            if (match) {
                clearTimeout(timer);
                res(match[1]);
            }
        });
        chrome.on('error', (err) => {
            clearTimeout(timer);
            rej(err);
        });
    });

    const port = wsUrl.match(/:(\d+)\//)[1];
    const newPageRes = await fetch(`http://127.0.0.1:${port}/json/new?file://${htmlPath}`, { method: 'PUT' });
    const target = await newPageRes.json();
    const pageWsUrl = target.webSocketDebuggerUrl;

    const ws = new WebSocket(pageWsUrl);
    let msgId = 1;
    const callbacks = new Map();

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        if (msg.id && callbacks.has(msg.id)) {
            callbacks.get(msg.id)(msg);
            callbacks.delete(msg.id);
        }
    };

    const send = (method, params = {}) => new Promise((res) => {
        const id = msgId++;
        callbacks.set(id, res);
        ws.send(JSON.stringify({ id, method, params }));
    });

    await new Promise((res) => ws.onopen = res);

    await send('Page.enable');
    await send('Runtime.enable');

    // Wait for web fonts (Poppins, JetBrains Mono, etc.) to finish downloading
    await send('Runtime.evaluate', {
        expression: 'document.fonts.ready.then(() => document.fonts.status)',
        awaitPromise: true
    });

    // Debounce for layout settle
    await new Promise((r) => setTimeout(r, 200));

    // Print to PDF with full theme background and 16:9 1080p landscape bounds
    const printRes = await send('Page.printToPDF', {
        printBackground: true,
        preferCSSPageSize: true,
        marginTop: 0,
        marginBottom: 0,
        marginLeft: 0,
        marginRight: 0
    });

    if (!printRes.result || !printRes.result.data) {
        throw new Error('Chrome printToPDF returned empty data');
    }

    const pdfBuffer = Buffer.from(printRes.result.data, 'base64');
    writeFileSync(outPdf, pdfBuffer);
    console.log(`Exported deck to: ${outPdf} (${pdfBuffer.length} bytes)`);

    try {
        ws.close();
        chrome.kill();
    } catch {}
}

main().catch((err) => {
    console.error('Export error:', err);
    process.exit(1);
});
