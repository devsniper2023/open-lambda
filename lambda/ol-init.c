#include <unistd.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/wait.h>

char **params;

// double fork to avoid zombies and exec the python server
void signal_handler() {
	if (fork() == 0) {
		if (fork() == 0) {
			execv(params[0], params);
		}
		exit(0);
	}

	wait(NULL);
	return;
}

int main(int argc, char *argv[]) {
	int k;

	params = (char**)malloc((3+argc-1)*sizeof(char*));
	params[0] = "/usr/bin/python";
	params[1] = "/server.py";
	for (k = 1; k < argc; k++) {
		params[k+1] = argv[k];
	}
	params[argc+1] = NULL;

	signal(SIGUSR1, signal_handler);
	pause(); // wait for SIGUSR1 and handle
	pause(); // sleep forever, we're init for the ns

	return 0;
}
