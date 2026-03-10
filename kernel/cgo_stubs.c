typedef long ssize_t;
typedef unsigned long size_t;
typedef unsigned long uint64_t;
typedef int int32_t;
typedef unsigned int uint32_t;
typedef int pid_t;
typedef unsigned int gid_t;
typedef unsigned int uid_t;
typedef int sigset_t;

void *stderr = (void*)0;
static int dummy_errno;
int *__errno_location(void) { return &dummy_errno; }

void abort(void) { for(;;) __asm__("hlt"); }
void kHalt(void) { for(;;) __asm__("hlt"); }
void *malloc(size_t n) { (void)n; return (void*)0; }
void free(void *p) { (void)p; }
int fprintf(void *f, const char *fmt, ...) { (void)f; (void)fmt; return 0; }
int fwrite(const void *p, size_t s, size_t n, void *f) { (void)p;(void)s;(void)n;(void)f; return 0; }
int fputc(int c, void *f) { (void)c;(void)f; return 0; }
int vfprintf(void *f, const char *fmt, __builtin_va_list ap) { (void)f;(void)fmt;(void)ap; return 0; }
char *strerror(int e) { (void)e; return (char*)""; }
void clearenv(void) {}
int setenv(const char *n, const char *v, int o) { (void)n;(void)v;(void)o; return 0; }
int unsetenv(const char *n) { (void)n; return 0; }
void *mmap(void *a, size_t l, int p, int f, int fd, long o) { (void)a;(void)l;(void)p;(void)f;(void)fd;(void)o; return (void*)-1; }
int munmap(void *a, size_t l) { (void)a;(void)l; return 0; }
int nanosleep(const void *req, void *rem) { (void)req;(void)rem; return 0; }

typedef void *pthread_t;
typedef void *pthread_mutex_t;
typedef void *pthread_cond_t;
typedef unsigned int pthread_key_t;
typedef void *pthread_attr_t;

int pthread_mutex_lock(pthread_mutex_t *m) { (void)m; return 0; }
int pthread_mutex_unlock(pthread_mutex_t *m) { (void)m; return 0; }
int pthread_cond_wait(pthread_cond_t *c, pthread_mutex_t *m) { (void)c;(void)m; return 0; }
int pthread_cond_broadcast(pthread_cond_t *c) { (void)c; return 0; }
int pthread_key_create(pthread_key_t *k, void (*d)(void*)) { (void)k;(void)d; return 0; }
int pthread_setspecific(pthread_key_t k, const void *v) { (void)k;(void)v; return 0; }
int pthread_create(pthread_t *t, const pthread_attr_t *a, void *(*f)(void*), void *arg) { (void)t;(void)a;(void)f;(void)arg; return 0; }
int pthread_attr_init(pthread_attr_t *a) { (void)a; return 0; }
int pthread_attr_setdetachstate(pthread_attr_t *a, int s) { (void)a;(void)s; return 0; }
int pthread_attr_getstacksize(const pthread_attr_t *a, size_t *s) { (void)a;(void)s; return 0; }
int pthread_attr_getstack(const pthread_attr_t *a, void **sp, size_t *ss) { (void)a;(void)sp;(void)ss; return 0; }
int pthread_attr_destroy(pthread_attr_t *a) { (void)a; return 0; }
pthread_t pthread_self(void) { return (pthread_t)0; }
int pthread_getattr_np(pthread_t t, pthread_attr_t *a) { (void)t;(void)a; return 0; }
int pthread_sigmask(int h, const sigset_t *n, sigset_t *o) { (void)h;(void)n;(void)o; return 0; }

int sigfillset(sigset_t *s) { (void)s; return 0; }
int sigemptyset(sigset_t *s) { (void)s; return 0; }
int sigaddset(sigset_t *s, int n) { (void)s;(void)n; return 0; }
int sigismember(const sigset_t *s, int n) { (void)s;(void)n; return 0; }
int sigaction(int n, const void *a, void *o) { (void)n;(void)a;(void)o; return 0; }

int setegid(gid_t g) { (void)g; return 0; }
int seteuid(uid_t u) { (void)u; return 0; }
int setgid(gid_t g) { (void)g; return 0; }
int setgroups(size_t n, const gid_t *l) { (void)n;(void)l; return 0; }
int setregid(gid_t r, gid_t e) { (void)r;(void)e; return 0; }
int setresgid(gid_t r, gid_t e, gid_t s) { (void)r;(void)e;(void)s; return 0; }
int setresuid(uid_t r, uid_t e, uid_t s) { (void)r;(void)e;(void)s; return 0; }
int setreuid(uid_t r, uid_t e) { (void)r;(void)e; return 0; }
int setuid(uid_t u) { (void)u; return 0; }
