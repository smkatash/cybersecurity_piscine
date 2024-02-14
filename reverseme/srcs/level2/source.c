#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void no(void) {
  puts("Nope.");
  exit(1);
}

void ok(void) {
  puts("Good job.");
  return;
}

int main() {
    char    input[24];
    char    pass[9];
    int     len_str;
    int     idx1;
    int     idx2;
    int     bval;
    int     num;
    char    buff[4];

    printf("Please enter key: ");
    if (scanf("%23s", input) != 1) {
        no();
    } 
    if (input[0] != '0') {
        no();
    }
    if (input[1] != '0') {
        no();
    }
    fflush(stdin);
    memset(pass, 0, 9);
    pass[0] = 'd';
    idx1 = 2;
    idx2 = 1;
    buff[3] = 0;
    while (1) {
        len_str = strlen(pass);
        bval = 0;
        if (len_str < 8) {
            len_str = strlen(input);
            bval = idx1 < len_str;
        }
        if (!bval) {
            break;
        }
        buff[0] = input[idx1];
        buff[1] = input[idx1 + 1];
        buff[2] = input[idx1 + 2];
        num = atoi(buff);
        pass[idx2] = (char)num;
        idx1 += 3;
        idx2 += 1;
    }
    pass[idx2] = 0;
    if (strcmp(pass, "delabere") == 0) {
        ok();
    } else {
        no();
    }
    return 0;
}
