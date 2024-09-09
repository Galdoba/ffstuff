Versioning:
    Major.Minor:Release Date

Changes done:
    v 0.1.0:2024-08-23 -- Prototype. 
    v 0.1.1:2024-08-26 -- Add direct processing.
    v 0.1.2:2024-08-26 -- Add bash generation.
    v 0.2.0:2024-08-28 -- Make config file not hardcoded.
    v 0.2.1:2024-08-28 -- Add auto renaming of source files based on $tv_series.xml
    v 0.2.2:2024-08-28 -- Add notifications
    v 0.2.3:2024-08-30 -- Reimplement sourcecollector's renaming module.
    v 0.3.0:2024-09-05 -- Implement Logging
    v 0.3.1:2024-09-06 -- bug fixing: renaming module

TODO:
    Major -- Develop report via *.ready file
    Major -- Export Task Formats to file
    Major -- Thread Separation: Error handling
    Major -- Thread Separation: Logging
    Major -- Movie Processing
    Major -- New Command: Health (check config/paths/etc...)
    Major -- New Mode: Prompt
    Major -- New Command: Menu (select other commands in prompt mode, set as default)
    
    Minor -- Unbind devtools/command library
    Minor -- Relocate to own Repo and create testing branch
    Minor -- Clean bash files
    Minor -- Move Source setup to Processing
    Minor -- Test: actions
    Minor -- Test: bashgen
    Minor -- Test: sourcefile
    Minor -- Test: targetfile
    Minor -- Test: bridge
    Minor -- Test: job
    Minor -- Test: media
    Minor -- Test: metainfo
    Minor -- Test: task
    Minor -- New Task: Archive
    Minor -- New Task: Store Stats 

Threads:
    Admin
    Error Handling
    DirectProcessing (?)
    TUI (?)
    Session Stats (?)
    Global Stats (?)