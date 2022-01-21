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

/**
 * Returns a string representation of the distinguished name,
 * roughly following the RFC 2253 Distinguished Names syntax.
 */
function formatDN(dn: DistinguishedName) {
  const attributes: [string, string[]][] = [
    ['C', []], // Country
    ['ST', []], // Province
    ['L', []], // Locality
    ['STREET', []], // StreetAddress
    ['POSTALCODE', []], // PostalCode
    ['O', []], // Organization
    ['OU', []], // OrganizationalUnit
    ['CN', []], // CommonName
    ['SERIALNUMBER', []], // SerialNumber
  ];

  for (const attribute of dn.attributes) {
    switch (attribute.oid) {
      case '2.5.4.6': // Country
        attributes[0][1].push(attribute.value);
        break;
      case '2.5.4.8': // Province
        attributes[1][1].push(attribute.value);
        break;
      case '2.5.4.7': // Locality
        attributes[2][1].push(attribute.value);
        break;
      case '2.5.4.9': // StreetAddress
        attributes[3][1].push(attribute.value);
        break;
      case '2.5.4.17': // PostalCode
        attributes[4][1].push(attribute.value);
        break;
      case '2.5.4.10': // Organization
        attributes[5][1].push(attribute.value);
        break;
      case '2.5.4.11': // OrganizationalUnit
        attributes[6][1].push(attribute.value);
        break;
      case '2.5.4.3': // Common Name (CN)
        attributes[7][1] = [attribute.value];
        break;
      case '2.5.4.5': // Serial Number
        attributes[8][1] = [attribute.value];
        break;
    }
  }

  const escape = (value: string) => {
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
  };

  return [...attributes]
    .reverse()
    .filter(([, values]) => values.length)
    .map(([key, values]) => values.map((value) => `${key}=${escape(value)}`).join('+'))
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
