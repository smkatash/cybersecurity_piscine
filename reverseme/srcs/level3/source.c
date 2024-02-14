#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void ___syscall_malloc(void) {
  puts("Nope.");
  exit(1);
}

void ____syscall_malloc(void) {
  puts("Good job.");
  return;
}


int main() {
    char input[31];
    char pass[9];
    char buff[4];
    int  idx1;
    int  idx2; 

    printf("Please enter key: ");
    if (scanf("%30s", input) != 1) {
        ___syscall_malloc();
    }
    if (input[0] != '4') {
        ___syscall_malloc();
    }
    if (input[1] != '2') {
        ___syscall_malloc();
    }
    fflush(stdin);
    memset(pass,0,9);
    pass[0] = '*';
    idx1 = 1;
    idx2 = 2;
    buff[3] = 0;
    int len_str;
    int bval;
    int num;
    while (1) {
        len_str = strlen(pass);
        bval = 0;
        if (len_str < 8) {
            len_str = strlen(input);
            bval = idx2 < len_str;
        }
        if (!bval) {
            break;
        }
        buff[0] = input[idx2];
        buff[1] = input[idx2 + 1];
        buff[2] = input[idx2 + 2];
        num = atoi(buff);
        pass[idx1] = (char)num;
        idx2 += 3;
        idx1 += 1;
    }
    pass[idx1] = 0;
    bval = strcmp(pass,"********");
      if (bval == -2) {
        ___syscall_malloc();
    } else if (bval == -1) {
        ___syscall_malloc();
    } else if (bval == 0) {
        ____syscall_malloc();
    } else if (bval == 1) {
        ___syscall_malloc();
    } else if (bval == 2) {
        ___syscall_malloc();
    } else if (bval == 3) {
        ___syscall_malloc();
    } else if (bval == 4) {
        ___syscall_malloc();
    } else if (bval == 5) {
        ___syscall_malloc();
    } else if (bval == 0x73) {
        ___syscall_malloc();
    } else {
        ___syscall_malloc();
    }
    return 0;
}