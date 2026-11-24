package rt

// Gate 4 extends the Gate 3 bootstrap inventory with the platform-thread,
// monitor and interrupt surface. The inventory is assembled in thread_java.go;
// this file documents the profile boundary. No application class may define
// anything in the bootstrap namespace, and earlier profiles never publish these
// members.
