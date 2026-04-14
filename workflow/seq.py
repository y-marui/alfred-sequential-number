import string


def convert_Alphabet(n):
    n = int(n) - 1
    if n < 0:
        raise ValueError("`n` must not be negative.")
    elif n == 0:
        return "A"
    else:
        s = ""
        while n > 0:
            if s == "":
                n, m = n // 26, n % 26
                s = string.ascii_uppercase[m] + s
            else:
                n, m = n // 27, n % 27
                s = string.ascii_uppercase[m - 1] + s
        return s


def convert_alphabet(n):
    n = int(n) - 1
    if n < 0:
        raise ValueError("`n` must not be negative.")
    elif n == 0:
        return "a"
    else:
        s = ""
        while n > 0:
            if s == "":
                n, m = n // 26, n % 26
                s = string.ascii_lowercase[m] + s
            else:
                n, m = n // 27, n % 27
                s = string.ascii_lowercase[m - 1] + s
        return s


class formatter(string.Formatter):
    def format_field(self, value, format_spec):
        if isinstance(value, int):
            if format_spec.endswith("a"):
                value = convert_alphabet(int(value))
                format_spec = format_spec[:-1]
            elif format_spec.endswith("A"):
                value = convert_Alphabet(int(value))
                format_spec = format_spec[:-1]
        return super().format_field(value, format_spec)


def load_range(text):
    if isinstance(text, int) or text.isdigit():
        return [1, int(text) + 1]
    else:
        return [int(s) + i for s, i in zip(text.split("-"), [0, 1])]


def multi_range(narr, *arg):
    if len(narr) == 1:
        for i in range(*load_range(narr[0])):
            yield *arg, i
    else:
        for i in range(*load_range(narr[0])):
            for arg2 in multi_range(narr[1:], *arg, i):
                yield arg2


def loader(pattern: str, narr: list):
    _flag = False
    _flag2 = False
    result = ""
    for c in pattern:
        if _flag:
            if c in ["b", "o", "d", "x", "X", "a", "A"]:
                result += "{:" + c + "}"
            else:
                raise ValueError("%b and %o, %d, %x, %X, %a, %A are only supported.")
            _flag = False
        elif _flag2:
            if c in ["a", "A"]:
                result += "{:" + c + "}"
            else:
                if c == "%":
                    _flag = True
                    result += "{:d}"
                else:
                    result += "{:d}" + c
            _flag2 = False
        elif c == "%":
            _flag = True
        elif c == "#":
            _flag2 = True
        else:
            result += c

    if _flag2:
        result += "{:d}"

    fmt = formatter()
    for arg in multi_range(narr):
        yield fmt.format(result, *arg)
