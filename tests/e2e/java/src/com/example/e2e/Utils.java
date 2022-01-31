package com.example.e2e;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.logging.Logger;

public class Utils {
  static {
    System.setProperty("java.util.logging.SimpleFormatter.format", "%1$tF %1$tT %5$s%n");
  }

  public static void log(Logger logger, String message) {
    logger.info(message);
  }

  public static void fatal(Logger logger, String message) {
    logger.severe(message);
    System.exit(1);
  }

  public static String[] getCredentials(
      String fabricRoot, String orgDomain, String username, String peerName) throws IOException {
    Path root = Path.of(fabricRoot, "test-network/organizations/peerOrganizations/", orgDomain);

    String[] paths =
        new String[] {
          "users/" + username + "/msp/signcerts/cert.pem",
          "users/" + username + "/msp/keystore/priv_sk",
          "peers/" + peerName + "/tls/ca.crt",
        };

    String[] credentials = new String[paths.length];
    for (int i = 0; i < paths.length; i++) {
      credentials[i] = Files.readString(Path.of(root.toString(), paths[i]));
    }

    return credentials;
  }
}
