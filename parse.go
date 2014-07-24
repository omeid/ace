package ace

import "strings"

// parseSource parses the source and returns the result.
func parseSource(src *source, opts *Options) (*result, error) {
	var rslt *result

	base, err := parseBytes(src.base.data, rslt, opts)
	if err != nil {
		return nil, err
	}

	inner, err := parseBytes(src.inner.data, rslt, opts)
	if err != nil {
		return nil, err
	}

	includes := make(map[string][]element)

	for _, f := range src.includes {
		includes[f.path], err = parseBytes(f.data, rslt, opts)
		if err != nil {
			return nil, err
		}
	}

	rslt = newResult(base, inner, includes)

	return rslt, nil
}

// parseBytes parses the byte data and returns the elements.
func parseBytes(data []byte, rslt *result, opts *Options) ([]element, error) {
	var elements []element

	lines := strings.Split(formatLF(string(data)), lf)

	i := 0
	l := len(lines)

	// Ignore the last empty line.
	if l > 0 && lines[l-1] == "" {
		l--
	}

	for i < l {
		// Fetch a line.
		ln := newLine(i+1, lines[i])
		i++

		// Ignore the empty line.
		if ln.isEmpty() {
			continue
		}

		if ln.isTopIndent() {
			e, err := newElement(ln, rslt, nil, opts)
			if err != nil {
				return nil, err
			}

			// Append child elements to the element.
			if err := appendChildren(e, lines, &i, l, rslt, opts); err != nil {
				return nil, err
			}

			elements = append(elements, e)
		}
	}

	return elements, nil
}

// appendChildren parses the lines and appends the children to the element.
func appendChildren(parent element, lines []string, i *int, l int, rslt *result, opts *Options) error {
	for *i < l {
		// Fetch a line.
		ln := newLine(*i+1, lines[*i])

		// Check if the line is a child of the parent.
		ok, err := ln.childOf(parent)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		child, err := newElement(ln, rslt, parent, opts)
		if err != nil {
			return err
		}

		parent.AppendChild(child)

		*i++

		if child.CanHaveChildren() {
			if err := appendChildren(child, lines, i, l, rslt, opts); err != nil {
				return err
			}
		}
	}

	return nil
}
