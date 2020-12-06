#include "sem.h"

int
semgetvalwrapper(int semid, int semnum) {
    return semctl(semid, semnum, GETVAL);
}

//https://github.com/ptt/pttbbs/blob/master/common/bbs/passwd.c#L75
int
semctlsetvalwrapper(int semid, int semnum, int val) {
   union semun s;
   s.val = val;
   return semctl(semid, semnum, SETVAL, s);
}

//https://github.com/ptt/pttbbs/blob/master/common/bbs/passwd.c#L87
int
semwaitwrapper(int semid, int semnum) {
    struct sembuf buf = {semnum, -1, SEM_UNDO};
    return semop(semid, &buf, 1);
}

//https://github.com/ptt/pttbbs/blob/master/common/bbs/passwd.c#L96
int
sempostwrapper(int semid, int semnum) {
    struct sembuf buf = {semnum, 1, SEM_UNDO};
    return semop(semid, &buf, 1);
}

int
semdestroywrapper(int semid, int semnum) {
    return semctl(semid, semnum, IPC_RMID);
}