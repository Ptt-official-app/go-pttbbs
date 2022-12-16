#ifndef __GO_PTTBBS_SEM_H__
#define __GO_PTTBBS_SEM_H__

// from: https://github.com/shubhros/drunkendeluge/blob/master/semaphore/semaphore.go

#include <sys/sem.h>
// /* https://comp.os.linux.development.system.narkive.com/rvJxp3Vb/union-variable-error-storage-size-isn-t-known */
#if defined(_SEM_SEMUN_UNDEFINED)
union semun {
  int val;                   /* value for SETVAL */
  struct semid_ds *buf;      /* buffer for IPC_STAT, IPC_SET */
  unsigned short int *array; /* array for GETALL, SETALL */
  struct seminfo *__buf;     /* buffer for IPC_INFO */
};
#endif

#ifndef SEM_R
#define SEM_R 0400
#endif

#ifndef SEM_A
#define SEM_A 0200
#endif

int semgetvalwrapper(int semid, int semnum);
int semctlsetvalwrapper(int semid, int semnum, int val);
int semdestroywrapper(int semid, int semnum);
int semwaitwrapper(int semid, int semnum);
int sempostwrapper(int semid, int semnum);

#endif  //__GO_PTTBBS_SEM_H__
