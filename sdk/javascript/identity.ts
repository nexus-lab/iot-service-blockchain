import { Certificate, DistinguishedName } from '@fidm/x509';

/**
 * Parse an X509 certificate from PEM string
 *
 * @param certificate PEM-encoded X509 certificate
 * @returns parsed X509 certificate
 */
export function parseCertificate(certificate: string) {
  return Certificate.fromPEM(Buffer.from(certificate, 'utf-8'));
}

const DN_ORDER = ['CN', 'SERIALNUMBER', 'C', 'L', 'ST', 'STREET', 'O', 'OU', 'POSTALCODE'];
const OID_MAP: { [key: string]: string } = {
  '2.5.4.3': 'CN',
  '2.5.4.5': 'SERIALNUMBER',
  '2.5.4.6': 'C',
  '2.5.4.7': 'L',
  '2.5.4.8': 'ST',
  '2.5.4.9': 'STREET',
  '2.5.4.10': 'O',
  '2.5.4.11': 'OU',
  '2.5.4.17': 'POSTALCODE',
};

function escapeDN(value: string) {
  return value
    .replace('\\', '\\\\')
    .replace(',', '\\,')
    .replace('+', '\\+')
    .replace('"', '\\"')
    .replace('<', '\\<')
    .replace('>', '\\>')
    .replace(';', '\\;')
    .replace(/^ /, '\\ ')
    .replace(/ $/, '\\ ')
    .replace(/^#/, '\\#');
}

function mapDN(dn: DistinguishedName) {
  const map = new Map<string, string[]>();

  for (const attribute of dn.attributes) {
    const type = OID_MAP[attribute.oid] ?? attribute.oid;

    // CN and SERIALNUMBER are single-value field
    if (!map.has(type) || type === 'CN' || type === 'SERIALNUMBER') {
      map.set(type, []);
    }

    map.get(type)?.push(attribute.value);
  }

  for (const values of map.values()) {
    values.sort();
  }

  return map;
}

/**
 * Returns a string representation of the distinguished name,
 * roughly following the RFC 2253 Distinguished Names syntax.
 * Distinguished Names are sorted by their OID name.
 */
function formatDN(dn: DistinguishedName) {
  const dnMap = mapDN(dn);

  return DN_ORDER.filter((type) => dnMap.has(type))
    .map((type) =>
      dnMap
        .get(type)
        ?.map((value) => `${type}=${escapeDN(value)}`)
        .join('+'),
    )
    .join(',');
}

/**
 * Get unique client ID from client certificate
 *
 * @param certificate client X509 certificate
 * @returns the unique client ID
 */
export function getClientId(certificate: Certificate) {
  const id = `x509::${formatDN(certificate.subject)}::${formatDN(certificate.issuer)}`;
  return Buffer.from(id, 'utf-8').toString('base64');
}
