import { spawn } from 'child_process';
import * as readline from 'readline';

async function runProcess(label: string, executable: string, ...args: string[]) {
  return new Promise<void>((resolve, reject) => {
    const process = spawn(executable, args);
    const stdout = readline.createInterface({ input: process.stdout });
    const stderr = readline.createInterface({ input: process.stderr });

    stdout.on('line', (line) => console.log(`${label} ${line}`));
    stderr.on('line', (line) => console.error(`${label} ${line}`));

    process.on('error', (err) => {
      stdout.close();
      stderr.close();
      reject(err);
    });

    process.on('exit', (code, signal) => {
      stdout.close();
      stderr.close();

      if (code === 0) {
        resolve();
      } else if (code !== null) {
        reject(new Error(`process exited with code ${code}`));
      } else {
        reject(new Error(`process exited with signal ${signal}`));
      }
    });
  });
}

async function main() {
  await Promise.all([
    runProcess('[eventhub]', 'node', 'tests/e2e/javascript/eventhub.js'),
    runProcess('[device]  ', 'node', 'tests/e2e/javascript/device.js'),
    runProcess('[app]     ', 'node', 'tests/e2e/javascript/app.js'),
  ]);
}

main();
