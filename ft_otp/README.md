
##### HMAC-SHA-1 calculation

HOTP(K,C) = Truncate(HMAC-SHA-1(K,C))
K and C represent the shared secret and counter value.

Symbol  Represents
-------------------------------------------------------------------
- C       8-byte counter value (8 x 8 => 64-bits), the moving factor.  This counter
           MUST be synchronized between the HOTP generator (client)
           and the HOTP validator (server).

- K       shared secret between client and server; each HOTP
           generator has a different and unique secret K.

#### Example of HOTP Computation for Digit = 6

The following code example describes the extraction of a dynamic
binary code given that hmac_result is a byte array with the HMAC-
SHA-1 result:
        int offset   =  hmac_result[19] & 0xf ;
        int bin_code = (hmac_result[offset]  & 0x7f) << 24
           | (hmac_result[offset+1] & 0xff) << 16
           | (hmac_result[offset+2] & 0xff) <<  8
           | (hmac_result[offset+3] & 0xff) ;

   SHA-1 HMAC Bytes (Example)

   -------------------------------------------------------------
   | Byte Number                                               |
   -------------------------------------------------------------
   |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|
   -------------------------------------------------------------
   | Byte Value                                                |
   -------------------------------------------------------------
   |1f|86|98|69|0e|02|ca|16|61|85|50|ef|7f|19|da|8e|94|5b|55|5a|
   -------------------------------***********----------------++|

   * The last byte (byte 19) has the hex value 0x5a.
   * The value of the lower 4 bits is 0xa (the offset value).
   * The offset value is byte 10 (0xa).
   * The value of the 4 bytes starting at byte 10 is 0x50ef7f19,
     which is the dynamic binary code DBC1.
   * The MSB of DBC1 is 0x50 so DBC2 = DBC1 = 0x50ef7f19 .
   * HOTP = DBC2 modulo 10^6 = 872921.

   We treat the dynamic binary code as a 31-bit, unsigned, big-endian
   integer; the first byte is masked with a 0x7f.

   We then take this number modulo 1,000,000 (10^6) to generate the 6-
   digit HOTP value 872921 decimal.
  
  
  Source [RFC 4226](https://www.rfc-editor.org/rfc/rfc4226)