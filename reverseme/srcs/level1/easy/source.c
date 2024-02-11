#include <stdio.h>
#include <string.h>
#include <stdint.h>

int32_t		main()
{
	char v4[114];

  	strcpy(v4, "__stack_check");
	printf("Please enter key: ");
	scanf("%s", &v4[14]);
	if (strcmp(&v4[14], v4) != 0)
	{
		printf("Nope.\n");
	}
	else
	{
		printf("Good job.\n");
	}
	return 0;
}