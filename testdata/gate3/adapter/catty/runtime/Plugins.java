package catty.runtime;

import java.io.IOException;
import java.net.URL;
import java.net.URLClassLoader;
import java.nio.file.Files;
import java.nio.file.Path;

/** Test-only OpenJDK 21 implementation of the Gate 3 API. */
public final class Plugins {
    private static final long LIMIT = 16L * 1024 * 1024;
    private Plugins() {}

    public static byte[] readJar(String path) throws IOException {
        if (path == null) {
            throw new NullPointerException("path");
        }
        Path p = Path.of(path);
        if (!Files.isRegularFile(p)) {
            throw new IOException("not a regular file: " + path);
        }
        long size = Files.size(p);
        if (size > LIMIT) {
            throw new PluginAdmissionException("LIMIT", "compressed input exceeds the limit");
        }
        return Files.readAllBytes(p);
    }

    public static ClassLoader register(byte[] jar, ClassLoader parent) {
        if (jar == null) {
            throw new NullPointerException("jar");
        }
        if (jar.length > LIMIT) {
            throw new PluginAdmissionException("LIMIT", "compressed input exceeds the limit");
        }
        byte[] copy = jar.clone();
        try {
            Path file = Files.createTempFile("catty-oracle-", ".jar");
            file.toFile().deleteOnExit();
            Files.write(file, copy);
            return new URLClassLoader(new URL[]{file.toUri().toURL()}, parent);
        } catch (IOException e) {
            throw new PluginAdmissionException("ARCHIVE", e.toString());
        }
    }
}
