import * as fs from 'fs';
import * as path from 'path';

function formatDate(date: Date) {
  const year = String(date.getFullYear()).padStart(4, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

export function log(message: string) {
  console.log(`${formatDate(new Date())} ${message}`);
}

export function fatal(message: string) {
  console.error(`${formatDate(new Date())} ${message}`);
  process.exit(1);
}

export function getCredentials(
  fabricRoot: string,
  orgDomain: string,
  username: string,
  peerName: string,
) {
  const root = path.join(fabricRoot, 'test-network/organizations/peerOrganizations/', orgDomain);
  const filepaths = [
    'users/' + username + '/msp/signcerts/cert.pem',
    'users/' + username + '/msp/keystore/priv_sk',
    'peers/' + peerName + '/tls/ca.crt',
  ];

  return filepaths.map((filepath) =>
    fs.readFileSync(path.join(root, filepath), { encoding: 'utf-8' }),
  );
}
