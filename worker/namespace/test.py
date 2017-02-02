import os, sys, ns, time
from subprocess import check_output

def main(pid):
    ls = ['ls', '/']
    ps = ['ps']
    sleep = ['sleep', '20']

    r = ns.forkenter(pid)

    # parent
    if r > 0:
        time.sleep(0.2)
        print("PARENT LS:\n%s\n" % check_output(ls).replace('\n', ' '))
        print("PARENT PS (pid=%s):\n%s" % (os.getpid(), check_output(ps)))

    # child
    if r < 0:
        time.sleep(0.1)
        print("CHILD LS:\n%s\n" % check_output(ls).replace('\n', ' '))
        print("CHILD PS (pid=%s):\n%s" % (os.getpid(), check_output(ps)))

    # grandchild
    if r == 0:
        print("GRANDCHILD LS:\n%s\n" % check_output(ls).replace('\n', ' '))
        print("GRANDCHILD PS (pid=%s):\n%s" % (os.getpid(), check_output(ps)))

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print('Usage: test.py <pid>')
        print('No PID provided, using 0.')
        main('0')
    else:
        main(sys.argv[1])
