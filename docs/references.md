# Directory standards references

Canonical documentation for the directory standards that toolpaths implements.

## XDG base directory specification

The primary standard for user-specific directories on Linux, FreeBSD, and OpenBSD. Defines environment variables and default paths for configuration, data, cache, state, and runtime directories.

- [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/latest/) - The current specification
- [basedir-spec wiki page](https://www.freedesktop.org/wiki/Specifications/basedir-spec/) - freedesktop.org project page

## macOS library directories

Apple's guidelines for where macOS apps should store configuration, data, caches, and logs. macOS uses Library subdirectories rather than XDG-style hidden files.

- [File System Programming Guide](https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/FileSystemProgrammingGuide/Introduction/Introduction.html) - Overview of file system concepts for macOS and iOS
- [macOS Library Directory Details](https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/FileSystemProgrammingGuide/MacOSXDirectories/MacOSXDirectories.html) - Reference for Library subdirectories and their purposes
- [Where to Put Application Files](https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPFileSystem/Articles/WhereToPutFiles.html) - Guidelines for app support files, caches, and temporary files
- [File System Basics](https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/FileSystemProgrammingGuide/FileSystemOverview/FileSystemOverview.html) - File system domains, security, and directory structure

## Windows known folders

Microsoft's system for identifying standard folders via unique identifiers. Replaces the older CSIDL constants. toolpaths uses `SHGetKnownFolderPath` to resolve paths.

- [Known Folders](https://learn.microsoft.com/en-us/windows/win32/shell/known-folders) - Overview of the Known Folders system
- [Working with Known Folders in Applications](https://learn.microsoft.com/en-us/windows/win32/shell/working-with-known-folders) - How to use Known Folder API functions
- [KNOWNFOLDERID](https://learn.microsoft.com/en-us/windows/win32/shell/knownfolderid) - Reference for all Known Folder identifiers
- [SHGetKnownFolderPath](https://learn.microsoft.com/en-us/windows/win32/api/shlobj_core/nf-shlobj_core-shgetknownfolderpath) - API function reference
- [CSIDL](https://learn.microsoft.com/en-us/windows/win32/shell/csidl) - Legacy constants (for reference)

## Filesystem hierarchy standard

The Linux Foundation standard for Unix-like system directory layouts. Provides background context for system-wide directory organization.

- [FHS 3.0](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/index.html) - Current specification (HTML)
- [FHS specifications](https://refspecs.linuxfoundation.org/fhs.shtml) - Landing page with all formats
- [FHS 3.0 PDF](https://refspecs.linuxfoundation.org/FHS_3.0/fhs-3.0.pdf) - PDF version

## Linux init system file hierarchy

The init system project's documentation on file system hierarchy. Incorporates and extends XDG and FHS for modern Linux systems.

- [file-hierarchy(7)](https://www.freedesktop.org/software/systemd/man/latest/file-hierarchy.html) - File hierarchy documentation

## FreeBSD and OpenBSD hierarchy

BSD-specific file system hierarchy documentation. While toolpaths follows XDG on BSD platforms, hier(7) provides context for system-wide paths.

- [FreeBSD hier(7)](https://man.freebsd.org/cgi/man.cgi?query=hier&sektion=7&format=html)
- [OpenBSD hier(7)](https://man.openbsd.org/hier)
- [Linux hier(7)](https://man7.org/linux/man-pages/man7/hier.7.html)
