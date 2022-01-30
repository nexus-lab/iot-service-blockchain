package org.nexus_lab.iot_service_blockchain.sdk;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThrows;

import java.security.cert.CertificateException;
import org.junit.Test;

public class IdentityTest {
  private static final String PUBLIC_KEY =
      "-----BEGIN PUBLIC KEY-----\n"
          + "MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty\n"
          + "WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3\n"
          + "CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x\n"
          + "qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/\n"
          + "yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY\n"
          + "nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/\n"
          + "6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q\n"
          + "TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/\n"
          + "a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9\n"
          + "PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ\n"
          + "yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk\n"
          + "AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==\n"
          + "-----END PUBLIC KEY-----";
  private static final String CERTIFICATE1 =
      "-----BEGIN CERTIFICATE-----\n"
          + "MIICmjCCAkCgAwIBAgIUd/uzCIgYnvr5IVrGgnVXIF/JvWMwCgYIKoZIzj0EAwIw\n"
          + "bDELMAkGA1UEBhMCVUsxEjAQBgNVBAgTCUhhbXBzaGlyZTEQMA4GA1UEBxMHSHVy\n"
          + "c2xleTEZMBcGA1UEChMQb3JnMi5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eub3Jn\n"
          + "Mi5leGFtcGxlLmNvbTAeFw0yMjAxMDYwMjE1MDBaFw0yMzAxMDYwMjIwMDBaMF0x\n"
          + "CzAJBgNVBAYTAlVTMRcwFQYDVQQIEw5Ob3J0aCBDYXJvbGluYTEUMBIGA1UEChML\n"
          + "SHlwZXJsZWRnZXIxDzANBgNVBAsTBmNsaWVudDEOMAwGA1UEAxMFdXNlcjEwWTAT\n"
          + "BgcqhkjOPQIBBggqhkjOPQMBBwNCAARe9edmNbHEx0pQJP3jfGgjtIDp0a/dmzR4\n"
          + "fi74zEQMKYz8E0nt/BTCGC8Uv9SRvBHI7biYW1k8WXfkCoPmPTjuo4HOMIHLMA4G\n"
          + "A1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQE+JosrrNPvToO\n"
          + "byzv7BkPFxp1QTAfBgNVHSMEGDAWgBQaaZoL4EglLGspr66g1a2vf83MvDARBgNV\n"
          + "HREECjAIggZIb21lUEMwWAYIKgMEBQYHCAEETHsiYXR0cnMiOnsiaGYuQWZmaWxp\n"
          + "YXRpb24iOiIiLCJoZi5FbnJvbGxtZW50SUQiOiJ1c2VyMSIsImhmLlR5cGUiOiJj\n"
          + "bGllbnQifX0wCgYIKoZIzj0EAwIDSAAwRQIhALxDGVIsgP3VxXMzrv+l0ijGgX4T\n"
          + "/AmTkI+tB0LZqzprAiAm3oeXhmFmxUXTnFXbumz7xelcodKByxXLHyAkucX/NA==\n"
          + "-----END CERTIFICATE-----";

  private static final String CERTIFICATE2 =
      "-----BEGIN CERTIFICATE-----\n"
          + "MIIF1DCCBLygAwIBAgIUD2bTOv4j4WIQsBum58McQeYbupMwDQYJKoZIhvcNAQEL\n"
          + "BQAwggFXMRQwEgYDVQQDDAtleGFtcGxlLmNvbTEZMBcGA1UEAwwQdGVzdC5leGFt\n"
          + "cGxlLmNvbTELMAkGA1UEBhMCVVMxCzAJBgNVBAYTAlVLMQswCQYDVQQIDAJUTjEU\n"
          + "MBIGA1UEBwwLQ2hhdHRhbm9vZ2ExJTAjBgNVBAkMHEFwdCAjMDAxLCAxMjM0IFNv\n"
          + "bWV3aGVyZSBTdC4xEjAQBgNVBAoMCU5leHVzIExhYjEXMBUGA1UECgwOVVQrQ2hh\n"
          + "dHRhbm9vZ2ExFjAUBgNVBAoMDVVUQ2hhdHRhbm9vZ2ExFzAVBgNVBAoMDlVUO0No\n"
          + "YXR0YW5vb2dhMQswCQYDVQQLDAJDUzEQMA4GA1UECwwHPEVNQ1M+OzEOMAwGA1UE\n"
          + "EQwFMzc0MDMxHjAcBgkqhkiG9w0BCQEWD2V4YW1wbGVAdXRjLmVkdTETMBEGA1UE\n"
          + "BRMKMTIzNDU2Kzc4OTAeFw0yMjAxMTkyMTQwMTZaFw0zMjAxMTcyMTQwMTZaMIIB\n"
          + "VzEUMBIGA1UEAwwLZXhhbXBsZS5jb20xGTAXBgNVBAMMEHRlc3QuZXhhbXBsZS5j\n"
          + "b20xCzAJBgNVBAYTAlVTMQswCQYDVQQGEwJVSzELMAkGA1UECAwCVE4xFDASBgNV\n"
          + "BAcMC0NoYXR0YW5vb2dhMSUwIwYDVQQJDBxBcHQgIzAwMSwgMTIzNCBTb21ld2hl\n"
          + "cmUgU3QuMRIwEAYDVQQKDAlOZXh1cyBMYWIxFzAVBgNVBAoMDlVUK0NoYXR0YW5v\n"
          + "b2dhMRYwFAYDVQQKDA1VVENoYXR0YW5vb2dhMRcwFQYDVQQKDA5VVDtDaGF0dGFu\n"
          + "b29nYTELMAkGA1UECwwCQ1MxEDAOBgNVBAsMBzxFTUNTPjsxDjAMBgNVBBEMBTM3\n"
          + "NDAzMR4wHAYJKoZIhvcNAQkBFg9leGFtcGxlQHV0Yy5lZHUxEzARBgNVBAUTCjEy\n"
          + "MzQ1Nis3ODkwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC4NCHY0n/9\n"
          + "YswvqB/Ht7fK3oyxVkq5yk2qgQ3Sb5Aw1uxcbvMmD55Ph3o419Z3SxMEvBCwykME\n"
          + "nLRSTkZVxqJNQKUsit3Nb9mXk5/6pFFycJHm/h8p199il2Ci44idj0rt4JRDCUK0\n"
          + "zF01DmbM8wA4fjyBEEkfeOiVwOjVYu9thb93+/O1MBMz4nsr0rHgG8JERCuYtJ2K\n"
          + "eGSLZJYECbSKtM5Hg18sf7VjG6864OhJYbI3VVCR8rNTd9TNl3gvgA9PqAUcisU6\n"
          + "8Mf9I6r8Ti0bNwDhTksfoM5aldBQbyHDwljmNJpJHqK5cJjeVHBX4TMtT3iNVWFh\n"
          + "fkrJFDB2WzAtAgMBAAGjgZMwgZAwHQYDVR0OBBYEFHoas4QTY+r6UH6WessVYR6P\n"
          + "tgHbMB8GA1UdIwQYMBaAFHoas4QTY+r6UH6WessVYR6PtgHbMA4GA1UdDwEB/wQE\n"
          + "AwIFoDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHAYDVR0RBBUw\n"
          + "E4ILZXhhbXBsZS5jb22HBAoAAAEwDQYJKoZIhvcNAQELBQADggEBAI51utUQFlJz\n"
          + "gviD+KFHTbCcfTp+aYqpAKUtzA/Z3H9TnGNQMA3ANFJGIYnazOTMdxw599mWsGnI\n"
          + "MGnn39bv52qgiH2u2/hFvagC6VLXtinTRM6gB/+E52eH7/OqPn3GRo5rAJl+q+/n\n"
          + "xntpUQy+PkbpMFIaMvIQvQFuKjXLHKNwb8Q4Yvu1uEdAzLYt4B9npKng8sNT9hbQ\n"
          + "t4iXMHAQMw5X/KYJXz32KQRpznjY6GZZixNo+IJxbFApEJ+ThGwBwsT22QYQo4X/\n"
          + "WIp4vGHDHVBKtXKT9PMjosNSOpWB6ffkyIinHQJa4RaFjZB01IZoCWx+5R58WclO\n"
          + "rJgTOS8SNUM=\n"
          + "-----END CERTIFICATE-----";

  private static final String CLIENT_ID1 =
      "eDUwOTo6Q049dXNlcjEsQz1VUyxTVD1Ob3J0aCBDYXJvbGluYSxPP"
          + "Uh5cGVybGVkZ2VyLE9VPWNsaWVudDo6Q049Y2Eub3JnMi5leGFtcGxlLmNvbSxDPVVLLEw9SHVyc2xleSxTV"
          + "D1IYW1wc2hpcmUsTz1vcmcyLmV4YW1wbGUuY29t";

  private static final String CLIENT_ID2 =
      "eDUwOTo6Q049dGVzdC5leGFtcGxlLmNvbSxTRVJJQUxOVU1CRVI9M"
          + "TIzNDU2XCs3ODksQz1VSytDPVVTLEw9Q2hhdHRhbm9vZ2EsU1Q9VE4sU1RSRUVUPUFwdCAjMDAxXCwgMTIzN"
          + "CBTb21ld2hlcmUgU3QuLE89TmV4dXMgTGFiK089VVRcK0NoYXR0YW5vb2dhK089VVRcO0NoYXR0YW5vb2dhK"
          + "089VVRDaGF0dGFub29nYSxPVT1cPEVNQ1NcPlw7K09VPUNTLFBPU1RBTENPREU9Mzc0MDM6OkNOPXRlc3QuZ"
          + "XhhbXBsZS5jb20sU0VSSUFMTlVNQkVSPTEyMzQ1NlwrNzg5LEM9VUsrQz1VUyxMPUNoYXR0YW5vb2dhLFNUP"
          + "VROLFNUUkVFVD1BcHQgIzAwMVwsIDEyMzQgU29tZXdoZXJlIFN0LixPPU5leHVzIExhYitPPVVUXCtDaGF0d"
          + "GFub29nYStPPVVUXDtDaGF0dGFub29nYStPPVVUQ2hhdHRhbm9vZ2EsT1U9XDxFTUNTXD5cOytPVT1DUyxQT"
          + "1NUQUxDT0RFPTM3NDAz";

  @Test
  public void testParseCertificate() throws Exception {
    assertThrows(
        CertificateException.class,
        () ->
            Identity.parseCertificate(
                "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----".getBytes()));
    assertThrows(
        CertificateException.class, () -> Identity.parseCertificate(PUBLIC_KEY.getBytes()));
    Identity.parseCertificate(CERTIFICATE1.getBytes());
    Identity.parseCertificate(CERTIFICATE2.getBytes());
  }

  @Test
  public void testGetClientId() throws Exception {
    assertEquals(
        CLIENT_ID1, Identity.getClientId(Identity.parseCertificate(CERTIFICATE1.getBytes())));
    assertEquals(
        CLIENT_ID2, Identity.getClientId(Identity.parseCertificate(CERTIFICATE2.getBytes())));
  }
}
