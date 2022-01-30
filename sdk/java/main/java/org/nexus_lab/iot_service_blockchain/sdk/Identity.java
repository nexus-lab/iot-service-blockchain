package org.nexus_lab.iot_service_blockchain.sdk;

import java.io.ByteArrayInputStream;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Base64;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import javax.naming.InvalidNameException;
import javax.naming.ldap.LdapName;
import javax.naming.ldap.Rdn;
import javax.security.auth.x500.X500Principal;

/** Untility functions for client identity */
public class Identity {
  private static final String[] DN_ORDER;
  private static final Map<String, String> OID_MAP = new HashMap<>();

  static {
    // RFC 2253 OID order
    DN_ORDER =
        new String[] {"CN", "SERIALNUMBER", "C", "L", "ST", "STREET", "O", "OU", "POSTALCODE"};

    OID_MAP.put("2.5.4.3", "CN");
    OID_MAP.put("2.5.4.5", "SERIALNUMBER");
    OID_MAP.put("2.5.4.6", "C");
    OID_MAP.put("2.5.4.7", "L");
    OID_MAP.put("2.5.4.8", "ST");
    OID_MAP.put("2.5.4.9", "STREET");
    OID_MAP.put("2.5.4.10", "O");
    OID_MAP.put("2.5.4.11", "OU");
    OID_MAP.put("2.5.4.17", "POSTALCODE");
  }

  /**
   * Parse an X509 certificate from PEM string
   *
   * @param certificate PEM-encoded X509 certificate
   * @return parsed X509 certificate
   * @throws CertificateException
   */
  public static X509Certificate parseCertificate(byte[] certificate) throws CertificateException {
    CertificateFactory factory = CertificateFactory.getInstance("X.509");
    return (X509Certificate) factory.generateCertificate(new ByteArrayInputStream(certificate));
  }

  private static String escapeDN(String dn) {
    return dn.replaceAll("\\\\", "\\\\\\\\")
        .replaceAll(",", "\\\\,")
        .replaceAll("\\+", "\\\\+")
        .replaceAll("\\\"", "\\\\\"")
        .replaceAll("<", "\\\\<")
        .replaceAll(">", "\\\\>")
        .replaceAll(";", "\\\\;")
        .replaceAll("^ ", "\\\\ ")
        .replaceAll(" $", "\\\\ ")
        .replaceAll("^#", "\\\\#");
  }

  private static Map<String, List<String>> mapDN(String dn) throws InvalidNameException {
    Map<String, List<String>> map = new HashMap<>();

    for (Rdn rdn : new LdapName(dn).getRdns()) {
      // CN and SERIALNUMBER are single-value field
      if (!map.containsKey(rdn.getType())
          || "CN".equals(rdn.getType())
          || "SERIALNUMBER".equals(rdn.getType())) {
        map.put(rdn.getType(), new ArrayList<>());
      }

      map.get(rdn.getType()).add(rdn.getValue().toString());
    }

    for (List<String> values : map.values()) {
      Collections.sort(values);
    }

    return map;
  }

  /**
   * Returns a string representation of the distinguished name, roughly following the RFC 2253
   * Distinguished Names syntax. Distinguished Names are sorted by their OID name.
   */
  private static String formatDN(String dn) throws InvalidNameException {
    Map<String, List<String>> dnMap = mapDN(dn);

    return Arrays.stream(DN_ORDER)
        .filter(type -> dnMap.containsKey(type))
        .map(
            type ->
                dnMap.get(type).stream()
                    .map(value -> type + "=" + escapeDN(value))
                    .collect(Collectors.joining("+")))
        .collect(Collectors.joining(","));
  }

  /**
   * Get unique client ID from client certificate
   *
   * @param certificate client X509 certificate
   * @return the unique client ID
   * @throws InvalidNameException
   */
  public static String getClientId(X509Certificate certificate) throws InvalidNameException {
    String subject = certificate.getSubjectX500Principal().getName(X500Principal.RFC2253, OID_MAP);
    String issuer = certificate.getIssuerX500Principal().getName(X500Principal.RFC2253, OID_MAP);
    String id = "x509::" + formatDN(subject) + "::" + formatDN(issuer);
    return Base64.getEncoder().encodeToString(id.getBytes());
  }
}
