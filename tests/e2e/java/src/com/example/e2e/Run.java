package com.example.e2e;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Callable;
import java.util.concurrent.Executors;
import java.util.concurrent.ThreadPoolExecutor;

public class Run {
  private static void runProcess(String label, String executable, String... args) {
    String[] cmd = new String[args.length + 1];
    cmd[0] = executable;
    System.arraycopy(args, 0, cmd, 1, args.length);

    try {
      ProcessBuilder builder = new ProcessBuilder(cmd);
      builder.environment().putAll(System.getenv());
      Process process = builder.redirectErrorStream(true).start();

      try (BufferedReader input =
          new BufferedReader(new InputStreamReader(process.getInputStream()))) {
        String line;
        while ((line = input.readLine()) != null) {
          System.out.println(label + " " + line);
        }
      }

      int exitCode = process.waitFor();
      if (exitCode != 0) {
        System.exit(exitCode);
      }
    } catch (Exception e) {
      e.printStackTrace();
      System.exit(1);
    }
  }

  public static void main(String[] args) throws Exception {
    ThreadPoolExecutor executor = (ThreadPoolExecutor) Executors.newFixedThreadPool(3);

    List<Callable<Object>> tasks = new ArrayList<>();
    tasks.add(
        Executors.callable(
            () ->
                runProcess(
                    "[eventhub]",
                    "mvn",
                    "exec:java",
                    "-Dexec.mainClass=com.example.e2e.EventHub")));
    tasks.add(
        Executors.callable(
            () ->
                runProcess(
                    "[device]  ", "mvn", "exec:java", "-Dexec.mainClass=com.example.e2e.App")));
    tasks.add(
        Executors.callable(
            () ->
                runProcess(
                    "[app]     ", "mvn", "exec:java", "-Dexec.mainClass=com.example.e2e.Device_")));
    executor.invokeAll(tasks);

    executor.shutdown();
  }
}
