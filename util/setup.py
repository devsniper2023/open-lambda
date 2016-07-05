#!/usr/bin/env python
import os, sys, random, string, argparse
from common import *

NGINX_EXAMPLE = 'docker run -d -p 80:80 -v %s:/usr/share/nginx/html:ro nginx'
def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--cluster', '-c', default='cluster')
    parser.add_argument('--appdir', '-d', default='')
    parser.add_argument('--appfile', '-f', default='')
    parser.add_argument('--scripts', '-s',  default='', nargs='*')
    args = parser.parse_args()

    appNames = os.listdir(os.path.join(SCRIPT_DIR, "..",  "applications"))
    if args.appdir not in appNames:
        print "That is not an application directory. Go to /applications."
        sys.exit()
    app_dir = os.path.join( SCRIPT_DIR, "..", "applications", args.appdir)
    app_files =  [z for z in os.listdir(app_dir) if os.path.isfile(os.path.join(app_dir, z))]   
    if args.appfile not in app_files:
        print "That file is not in this directory"
        sys.exit()
    for a in args.scripts:
        if a not in app_files:
            print a + " is not in the "+  args.appdir + " application directory."
            sys.exit()

    #print args.scripts
    #sys.exit()

    lambdaFn = args.appfile
    app_name = ''.join(random.choice(string.ascii_lowercase) for _ in range(12))
    static_dir = os.path.join(app_dir, 'static')
    root_dir = os.path.join(app_dir, '..', '..')
    cluster_dir = os.path.join(root_dir, 'util', args.cluster)
    builder_dir = os.path.join(root_dir, 'lambda-generator')
    if not os.path.exists(cluster_dir):
        return 'cluster not running'

    # build image
    print '='*40
    print 'Building image'
    builder = os.path.join(builder_dir, 'builder.py')
    builder = builder + ' -n %s -l %s' %(app_name, os.path.join(app_dir, lambdaFn))
    if 'lambda-config.json' in app_files:
        builder = builder + ' -c %s' %(os.path.join(app_dir, 'lambda-config.json'))
    if 'environment.json' in app_files:
        builder = builder + ' -e %s' %(os.path.join(app_dir, 'environment.json'))

    run(builder, True)

    # push image
    print '='*40
    print 'Pushing image'
    registry = rdjs(os.path.join(cluster_dir, 'registry.json'))
    img = 'localhost:%s/%s' % (registry['host_port'], app_name)
    run('docker tag -f %s %s' % (app_name, img), True)
    run('docker push ' + img, True)

    # setup config
    balancer = rdjs(os.path.join(cluster_dir, 'loadbalancer-1.json'))
    config_file = os.path.join(static_dir, 'config.json')
    url = ("http://%s:%s/runLambda/%s" %
           (balancer['host_ip'], balancer['host_port'], app_name))
    wrjs(config_file, {'url': url})

    #run additional scripts, if there are any
    print '='*40
    if args.scripts == '':
        print "No additional scripts to run"
    else:
        print "Running additional scripts"
        for scr in args.scripts:
            spath = os.path.join(app_dir, scr)
            spath = "python " + spath
            print '='*40
            print "running " + scr
            run(spath, True)


            
    # directions
    print '='*40
    print 'Consider serving the app with nginx as follows:'
    print NGINX_EXAMPLE % static_dir
    return None

if __name__ == '__main__':
    rv = main()
    if rv != None:
        print 'ERROR: ' + rv
        sys.exit(1)
    sys.exit(0)
