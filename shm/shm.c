#include "shm.h"

char *_BCACHEPTR;

void
readwrapper(void *outptr, void *shmaddr, int offset, unsigned long n) {
  unsigned char *src = (unsigned char *)shmaddr + offset;
  memcpy(outptr, src, n);
}

void
writewrapper(void *shmaddr, int offset, void *inptr, unsigned long n) {
  unsigned char *dst = (unsigned char *)shmaddr + offset;
  memcpy(dst, inptr, n);
}

void
incuint32wrapper(void *shmaddr, int offset) {
  unsigned char *dst_b = (unsigned char *)shmaddr + offset;
  unsigned int *dst = (unsigned int *)dst_b;
  (*dst)++;
}

void
set_or_uint32wrapper(void *shmaddr, int offset, unsigned int flag) {
    unsigned char *c_dst = (unsigned char *)shmaddr + offset;
    unsigned int *dst = (unsigned int *)c_dst;
    *dst |= flag;
}

void
innerset_int32wrapper(void *shmaddr, int offsetSrc, int offsetDst) {
    unsigned char *c_src = (unsigned char *)shmaddr + offsetSrc;
    unsigned char *c_dst = (unsigned char *)shmaddr + offsetDst;
    int *src = (int *)c_src;
    int *dst = (int *)c_dst;
    *dst = *src;
}

int
memcmpwrapper(void *shmaddr, int offset, unsigned long n, void *cmpaddr) {
  unsigned char *cmp1 = (unsigned char *)shmaddr + offset;
  return memcmp(cmp1, cmpaddr, n);
}

void
memsetwrapper(void *shmaddr, int offset, unsigned char c, unsigned long n) {
  unsigned char *dst = (unsigned char *)shmaddr + offset;
  memset(dst, c, n);
}

void
set_bcacheptr(void *shmaddr, int offset) {
    _BCACHEPTR = (char *)shmaddr + offset;
}

char *
_bcache(int idx) {
    return _BCACHEPTR + SIZE_BOARD_HEADER*idx;
}

void
qsort_cmpboardname_wrapper(void *shmaddr, int offset, unsigned long n) {
    //the bsortedaddr is already in board-name
    unsigned char *bsortedaddr = (unsigned char *)shmaddr + offset;

    //init bsorted
    int i = 0;
    int *pbsortedaddr = (int *)bsortedaddr;
    for (i = 0, pbsortedaddr = (int *)bsortedaddr; i < n; i++, pbsortedaddr++) {
        *pbsortedaddr = i;
    }

    //qsort
    qsort(bsortedaddr, n, sizeof(int), cmpboardname);
}

void
qsort_cmpboardclass_wrapper(void *shmaddr, int offset, unsigned long n) {
    //the bsortedaddr is already in board-class
    unsigned char *bsortedaddr = (unsigned char *)shmaddr + offset;

    //init bsorted
    int i = 0;
    int *pbsortedaddr = (int *)bsortedaddr;
    for (i = 0, pbsortedaddr = (int *)bsortedaddr; i < n; i++, pbsortedaddr++) {
        *pbsortedaddr = i;
    }

    //qsort
    qsort(bsortedaddr, n, sizeof(int), cmpboardclass);
}

/**
 * qsort comparison function - 照板名排序
 */
int
cmpboardname(const void * i, const void * j) {
    int int_i = *(int *)i;
    int int_j = *(int *)j;
    char *bcache_i = _bcache(int_i);
    char *bcache_j = _bcache(int_j);
    char *bcache_cmp_i = bcache_i + OFFSET_BOARD_HEADER_BRDNAME;
    char *bcache_cmp_j = bcache_j + OFFSET_BOARD_HEADER_BRDNAME;

    return strncasecmp(bcache_cmp_i, bcache_cmp_j, SIZE_BOARD_HEADER_BRDNAME);
}

/**
 * qsort comparison function - 先照群組排序、同一個群組內依板名排
 */
int
cmpboardclass(const void * i, const void * j) {
    char *bcache_i = _bcache(*(int *)i);
    char *bcache_j = _bcache(*(int *)j);
    char *bcache_cmp_i = bcache_i + OFFSET_BOARD_HEADER_TITLE;
    char *bcache_cmp_j = bcache_j + OFFSET_BOARD_HEADER_TITLE;
    int cmp;

    cmp=strncmp(bcache_cmp_i, bcache_cmp_j, 4);
    if(cmp!=0) return cmp;

    bcache_cmp_i = bcache_i + OFFSET_BOARD_HEADER_BRDNAME;
    bcache_cmp_j = bcache_j + OFFSET_BOARD_HEADER_BRDNAME;

    return strncasecmp(bcache_cmp_i, bcache_cmp_j, SIZE_BOARD_HEADER_BRDNAME);
}
