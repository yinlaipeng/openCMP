#!/bin/bash
#
# Spec Queue Helpers for Ralph Loop
#
# Functions to count and find incomplete specs
#

# Count all root-level spec directories in specs/ folder
# A root spec is a directory containing a spec.md file
count_root_specs() {
    local specs_dir="$1"
    local count=0

    if [ ! -d "$specs_dir" ]; then
        echo 0
        return
    fi

    # Count directories that contain spec.md
    for dir in "$specs_dir"/*/; do
        [ -d "$dir" ] || continue
        if [ -f "$dir/spec.md" ]; then
            count=$((count + 1))
        fi
    done

    echo $count
}

# Count incomplete root specs (those without "Status: COMPLETE")
count_incomplete_root_specs() {
    local specs_dir="$1"
    local count=0

    if [ ! -d "$specs_dir" ]; then
        echo 0
        return
    fi

    for dir in "$specs_dir"/*/; do
        [ -d "$dir" ] || continue
        local spec_file="$dir/spec.md"
        [ -f "$spec_file" ] || continue

        # Check if spec is incomplete (no "Status: COMPLETE")
        if ! grep -q "Status: COMPLETE" "$spec_file" 2>/dev/null; then
            count=$((count + 1))
        fi
    done

    echo $count
}

# Get the first incomplete root spec by priority (lowest number = highest priority)
# Returns the spec directory path relative to specs/
get_first_incomplete_root_spec() {
    local specs_dir="$1"
    local first_spec=""
    local lowest_num=999999

    if [ ! -d "$specs_dir" ]; then
        echo ""
        return
    fi

    for dir in "$specs_dir"/*/; do
        [ -d "$dir" ] || continue
        local spec_file="$dir/spec.md"
        [ -f "$spec_file" ] || continue

        # Skip if complete
        if grep -q "Status: COMPLETE" "$spec_file" 2>/dev/null; then
            continue
        fi

        # Extract priority number from directory name (e.g., "001-xxx" -> 001)
        local dirname=$(basename "$dir")
        local num=$(echo "$dirname" | grep -oE '^[0-9]+' | sed 's/^0*//' || echo "999999")
        num=${num:-999999}

        if [ "$num" -lt "$lowest_num" ]; then
            lowest_num=$num
            first_spec="$dirname"
        fi
    done

    echo "$first_spec"
}

# Get the spec file path for a given spec directory
get_spec_file() {
    local specs_dir="$1"
    local spec_dirname="$2"

    echo "$specs_dir/$spec_dirname/spec.md"
}

# Check if a spec is complete
is_spec_complete() {
    local spec_file="$1"

    if [ ! -f "$spec_file" ]; then
        echo "false"
        return
    fi

    if grep -q "Status: COMPLETE" "$spec_file" 2>/dev/null; then
        echo "true"
    else
        echo "false"
    fi
}