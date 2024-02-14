#include <stdio.h>
#include <string.h>

int		main()
{
	char input[100];
	char *pass = "__stack_check";

	printf("Please enter key: ");
	scanf("%s", input);
	if (strcmp(input, pass) == 0) {
		printf("Good job.\n");
	} else {
		printf("Nope.\n");
	}
	return 0;
}