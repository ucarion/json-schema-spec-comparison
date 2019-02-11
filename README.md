# json-schema-spec-comparison

This repo documents the known differences between JSON Schema, as implemented by
[`json-schema-spec/json-schema-go`][json-schema-spec], and IETF I-D
`draft-handrews-01-jsonschema` ("draft-07"), as described by
[`json-schema-org/JSON-Schema-Test-Suite`][json-schema-org].

[json-schema-spec]: https://github.com/json-schema-spec/json-schema-go
[json-schema-org]: https://github.com/json-schema-org/JSON-Schema-Test-Suite

## Findings

This section describes the differences found in this comparison. In executive
summary, the differences are captured by the following statement:

> `json-schema-go` strives to avoid insecure, poorly-defined, or confusing
   behavior.

Therefore, `json-schema-go` avoids:

1. Auto-fetching schemas (insecure),
2. Auto-assigning IDs to schemas (insecure),
3. Changing base URIs (poorly-defined),
4. Having `$ref` disable sibling keywords (confusing),
5. The `format` keyword (poorly-defined),
6. The `contentMediaType` and `contentEncoding` keywords (poorly-defined),
7. Auto-promoting numbers to Bignum (confusing),
8. Emulating ECMAScript regexes (insecure),

In practice, it is dubious whether many real-life use-cases are very much
affected by these changes. This decision comes down to the individual's
use-case.

### Automatically Fetching Schemas

The JSON Schema Test Suite does not work as-is with `json-schema-go`. This is
for two reasons:

1. The JSON Schema Test Suite presumes that validators automatically fetch
   schemas when they are referenced by URL. `json-schema-go` does not do this
   for security reasons.
2. The JSON Schema Test suite presumes that when a validator fetches a schema by
   its URL, it automatically assigns that schema an ID based on the URL it
   fetched it from. Because `json-schema-go` does *not* automatically fetch
   schemas, this presumption is unsatisfiable.

This is a pretty minor difference in practice. It's a bad idea to host JSON
Schemas online, but not give them an `$id`, so few real-world use-cases will be
affected by this change.

The workaround, for the purposes of this comparison, is pretty simple -- the
schemas that are meant to be loaded over HTTP are simply pre-loaded directly
from the filesystem, and those the schemas were modified by hand to include the
appropriate `$id`.

### Changing Base URI

The JSON Schema Test Suite requires that validators (and schema authors) support
chaging the "base URI" of references from *any* sub-schema. This is a perennial
source of troubles for JSON Schema, and is the subject of much confusion.

This ticket, in particular, addresses a potentially unsolvable problem with
`$id` as specified by JSON Schema Test Suite -- that it is impossible to
determine when `$id` should change base URIs or not:

https://github.com/json-schema-org/json-schema-spec/issues/687

In an effort to cut the Gordian knot, `json-schema-go` ignores the `$id` keyword
outside of root schemas. Root schemas are the only place where `$id` is
well-defined and unsurprising.

### `$ref` Disabling Sibling Keywords

The JSON Schema Test Suite requires that the `$ref` keyword disable all sibling
keywords. `json-schema-go` does *not* adhere to this, because it's quite
confusing and without clear benefits to the user.

In other words, `$ref` in `json-schema-go` works like all other keywords. It's
implicitly ANDed will all of its siblings.

### MutipleOf Precision

The tests in the suite for `multipleOf` are a bit bizarre. `json-schema-go` uses
the following test to determine if *N* is a multiple of *M*:

> `| N % M | < Epsilon`

There doesn't exist a single `Epsilon` value that satisfies the suite. In
practice, this is rarely going to be an issue, because few people care to
distinguish between nearly-integral and truly-integral values in a way that a
single `Epsilon` cannot satisfy.

### Optional Test Cases

The test suite includes optional test cases. Of those optional test cases, the
following are not supported:

* `json-schema-go` does not support the `format` keyword. Instead, users should
  use the `pattern` keyword, which is far more portable.
* `json-schema-go` does not support the `contentMediaType` and `contentEncoding`
  keywords. This is because few people use it in practice, and because
  supporting it requires supporting *all* media types and encodings, an
  impossible task.
* `json-schema-go` does not support automatically using "Bignum" instead of
  floating-point numbers. This is by design, and both for for performance
  reasons and because it's surprising behavior. JSON is formally defined in
  terms of IEEE Floating-Point numbers, and `json-schema-go` sticks to this
  definition.
* `json-schema-go` does not emulate the ECMAScript Regex syntax because this is
  a Sisyphean task, and also opens the door to ReDoS, because ECMAScript regular
  expressions do not have the linear-time guarantees that Golang's do.

## Raw Output

Here's the raw output of running the tests. See above for discussion -- all
failures are accounted for and intentional:

```
go test -json ./... | jq -r 'select(.Action=="fail") | .Test | select(.!=null)'
```

```
TestSpec/tests/draft7/multipleOf.json/by_small_number/0.00751_is_not_multiple_of_0.0001
TestSpec/tests/draft7/multipleOf.json/by_small_number
TestSpec/tests/draft7/multipleOf.json
TestSpec/tests/draft7/optional/bignum.json/float_comparison_with_high_precision/comparison_works_for_high_numbers
TestSpec/tests/draft7/optional/bignum.json/float_comparison_with_high_precision
TestSpec/tests/draft7/optional/bignum.json/float_comparison_with_high_precision_on_negative_numbers/comparison_works_for_very_negative_numbers
TestSpec/tests/draft7/optional/bignum.json/float_comparison_with_high_precision_on_negative_numbers
TestSpec/tests/draft7/optional/bignum.json
TestSpec/tests/draft7/optional/content.json/validation_of_string-encoded_content_based_on_media_type/an_invalid_JSON_document
TestSpec/tests/draft7/optional/content.json/validation_of_string-encoded_content_based_on_media_type
TestSpec/tests/draft7/optional/content.json/validation_of_binary_string-encoding/an_invalid_base64_string_(%_is_not_a_valid_character)
TestSpec/tests/draft7/optional/content.json/validation_of_binary_string-encoding
TestSpec/tests/draft7/optional/content.json/validation_of_binary-encoded_media_type_documents/a_validly-encoded_invalid_JSON_document
TestSpec/tests/draft7/optional/content.json/validation_of_binary-encoded_media_type_documents/an_invalid_base64_string_that_is_valid_JSON
TestSpec/tests/draft7/optional/content.json/validation_of_binary-encoded_media_type_documents
TestSpec/tests/draft7/optional/content.json
TestSpec/tests/draft7/optional/ecmascript-regex.json/ECMA_262_regex_non-compliance/ECMA_262_has_no_support_for_\Z_anchor_from_.NET
TestSpec/tests/draft7/optional/ecmascript-regex.json/ECMA_262_regex_non-compliance
TestSpec/tests/draft7/optional/ecmascript-regex.json
TestSpec/tests/draft7/optional/format/date-time.json/validation_of_date-time_strings/a_invalid_day_in_date-time_string
TestSpec/tests/draft7/optional/format/date-time.json/validation_of_date-time_strings/an_invalid_offset_in_date-time_string
TestSpec/tests/draft7/optional/format/date-time.json/validation_of_date-time_strings/an_invalid_date-time_string
TestSpec/tests/draft7/optional/format/date-time.json/validation_of_date-time_strings/only_RFC3339_not_all_of_ISO_8601_are_valid
TestSpec/tests/draft7/optional/format/date-time.json/validation_of_date-time_strings
TestSpec/tests/draft7/optional/format/date-time.json
TestSpec/tests/draft7/optional/format/date.json/validation_of_date_strings/an_invalid_date-time_string
TestSpec/tests/draft7/optional/format/date.json/validation_of_date_strings/only_RFC3339_not_all_of_ISO_8601_are_valid
TestSpec/tests/draft7/optional/format/date.json/validation_of_date_strings
TestSpec/tests/draft7/optional/format/date.json
TestSpec/tests/draft7/optional/format/email.json/validation_of_e-mail_addresses/an_invalid_e-mail_address
TestSpec/tests/draft7/optional/format/email.json/validation_of_e-mail_addresses
TestSpec/tests/draft7/optional/format/email.json
TestSpec/tests/draft7/optional/format/hostname.json/validation_of_host_names/a_host_name_starting_with_an_illegal_character
TestSpec/tests/draft7/optional/format/hostname.json/validation_of_host_names/a_host_name_containing_illegal_characters
TestSpec/tests/draft7/optional/format/hostname.json/validation_of_host_names/a_host_name_with_a_component_too_long
TestSpec/tests/draft7/optional/format/hostname.json/validation_of_host_names
TestSpec/tests/draft7/optional/format/hostname.json
TestSpec/tests/draft7/optional/format/idn-email.json/validation_of_an_internationalized_e-mail_addresses/an_invalid_idn_e-mail_address
TestSpec/tests/draft7/optional/format/idn-email.json/validation_of_an_internationalized_e-mail_addresses
TestSpec/tests/draft7/optional/format/idn-email.json
TestSpec/tests/draft7/optional/format/idn-hostname.json/validation_of_internationalized_host_names/illegal_first_char_U+302E_Hangul_single_dot_tone_mark
TestSpec/tests/draft7/optional/format/idn-hostname.json/validation_of_internationalized_host_names/contains_illegal_char_U+302E_Hangul_single_dot_tone_mark
TestSpec/tests/draft7/optional/format/idn-hostname.json/validation_of_internationalized_host_names/a_host_name_with_a_component_too_long
TestSpec/tests/draft7/optional/format/idn-hostname.json/validation_of_internationalized_host_names
TestSpec/tests/draft7/optional/format/idn-hostname.json
TestSpec/tests/draft7/optional/format/ipv4.json/validation_of_IP_addresses/an_IP_address_with_too_many_components
TestSpec/tests/draft7/optional/format/ipv4.json/validation_of_IP_addresses/an_IP_address_with_out-of-range_values
TestSpec/tests/draft7/optional/format/ipv4.json/validation_of_IP_addresses/an_IP_address_without_4_components
TestSpec/tests/draft7/optional/format/ipv4.json/validation_of_IP_addresses/an_IP_address_as_an_integer
TestSpec/tests/draft7/optional/format/ipv4.json/validation_of_IP_addresses
TestSpec/tests/draft7/optional/format/ipv4.json
TestSpec/tests/draft7/optional/format/ipv6.json/validation_of_IPv6_addresses/an_IPv6_address_with_out-of-range_values
TestSpec/tests/draft7/optional/format/ipv6.json/validation_of_IPv6_addresses/an_IPv6_address_with_too_many_components
TestSpec/tests/draft7/optional/format/ipv6.json/validation_of_IPv6_addresses/an_IPv6_address_containing_illegal_characters
TestSpec/tests/draft7/optional/format/ipv6.json/validation_of_IPv6_addresses
TestSpec/tests/draft7/optional/format/ipv6.json
TestSpec/tests/draft7/optional/format/iri-reference.json/validation_of_IRI_References/an_invalid_IRI_Reference
TestSpec/tests/draft7/optional/format/iri-reference.json/validation_of_IRI_References/an_invalid_IRI_fragment
TestSpec/tests/draft7/optional/format/iri-reference.json/validation_of_IRI_References
TestSpec/tests/draft7/optional/format/iri-reference.json
TestSpec/tests/draft7/optional/format/iri.json/validation_of_IRIs/an_invalid_IRI_based_on_IPv6
TestSpec/tests/draft7/optional/format/iri.json/validation_of_IRIs/an_invalid_relative_IRI_Reference
TestSpec/tests/draft7/optional/format/iri.json/validation_of_IRIs/an_invalid_IRI
TestSpec/tests/draft7/optional/format/iri.json/validation_of_IRIs/an_invalid_IRI_though_valid_IRI_reference
TestSpec/tests/draft7/optional/format/iri.json/validation_of_IRIs
TestSpec/tests/draft7/optional/format/iri.json
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(~_not_escaped)
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(URI_Fragment_Identifier)_#1
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(URI_Fragment_Identifier)_#2
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(URI_Fragment_Identifier)_#3
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(some_escaped,_but_not_all)_#1
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(some_escaped,_but_not_all)_#2
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(wrong_escape_character)_#1
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(wrong_escape_character)_#2
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(multiple_characters_not_escaped)
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(isn't_empty_nor_starts_with_/)_#1
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(isn't_empty_nor_starts_with_/)_#2
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)/not_a_valid_JSON-pointer_(isn't_empty_nor_starts_with_/)_#3
TestSpec/tests/draft7/optional/format/json-pointer.json/validation_of_JSON-pointers_(JSON_String_Representation)
TestSpec/tests/draft7/optional/format/json-pointer.json
TestSpec/tests/draft7/optional/format/regex.json/validation_of_regular_expressions/a_regular_expression_with_unclosed_parens_is_invalid
TestSpec/tests/draft7/optional/format/regex.json/validation_of_regular_expressions
TestSpec/tests/draft7/optional/format/regex.json
TestSpec/tests/draft7/optional/format/relative-json-pointer.json/validation_of_Relative_JSON_Pointers_(RJP)/an_invalid_RJP_that_is_a_valid_JSON_Pointer
TestSpec/tests/draft7/optional/format/relative-json-pointer.json/validation_of_Relative_JSON_Pointers_(RJP)
TestSpec/tests/draft7/optional/format/relative-json-pointer.json
TestSpec/tests/draft7/optional/format/time.json/validation_of_time_strings/an_invalid_time_string
TestSpec/tests/draft7/optional/format/time.json/validation_of_time_strings/only_RFC3339_not_all_of_ISO_8601_are_valid
TestSpec/tests/draft7/optional/format/time.json/validation_of_time_strings
TestSpec/tests/draft7/optional/format/time.json
TestSpec/tests/draft7/optional/format/uri-reference.json/validation_of_URI_References/an_invalid_URI_Reference
TestSpec/tests/draft7/optional/format/uri-reference.json/validation_of_URI_References/an_invalid_URI_fragment
TestSpec/tests/draft7/optional/format/uri-reference.json/validation_of_URI_References
TestSpec/tests/draft7/optional/format/uri-reference.json
TestSpec/tests/draft7/optional/format/uri-template.json/format:_uri-template/an_invalid_uri-template
TestSpec/tests/draft7/optional/format/uri-template.json/format:_uri-template
TestSpec/tests/draft7/optional/format/uri-template.json
TestSpec/tests/draft7/optional/format/uri.json/validation_of_URIs/an_invalid_protocol-relative_URI_Reference
TestSpec/tests/draft7/optional/format/uri.json/validation_of_URIs/an_invalid_relative_URI_Reference
TestSpec/tests/draft7/optional/format/uri.json/validation_of_URIs/an_invalid_URI
TestSpec/tests/draft7/optional/format/uri.json/validation_of_URIs/an_invalid_URI_though_valid_URI_reference
TestSpec/tests/draft7/optional/format/uri.json/validation_of_URIs/an_invalid_URI_with_spaces
TestSpec/tests/draft7/optional/format/uri.json/validation_of_URIs/an_invalid_URI_with_spaces_and_missing_scheme
TestSpec/tests/draft7/optional/format/uri.json/validation_of_URIs
TestSpec/tests/draft7/optional/format/uri.json
TestSpec/tests/draft7/ref.json/ref_overrides_any_sibling_keywords/ref_valid,_maxItems_ignored
TestSpec/tests/draft7/ref.json/ref_overrides_any_sibling_keywords
TestSpec/tests/draft7/ref.json/Recursive_references_between_schemas
TestSpec/tests/draft7/ref.json
TestSpec/tests/draft7/refRemote.json/base_URI_change
TestSpec/tests/draft7/refRemote.json/base_URI_change_-_change_folder
TestSpec/tests/draft7/refRemote.json/base_URI_change_-_change_folder_in_subschema
TestSpec/tests/draft7/refRemote.json
TestSpec
```
