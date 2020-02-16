/*
Package working provides types for file system operations within a directory.

   d := new(working.Directory)
   d.Touch("somefile")
   body, err := d.Load("config")

or working with temporary directories

   d := new(working.Directory)
   d.Temporary()
   defer d.RemoveAll()
   // .. do stuff

*/
package working
