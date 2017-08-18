typedef struct DriveStats DriveStats;

enum {
	NDRIVE = 16,
	NAMELEN = 31
};

struct DriveStats {
	char name[NAMELEN+1];
	int64_t size;
	int64_t blocksize;

	int64_t read;
	int64_t written;
	int64_t nread;
	int64_t nwrite;
	int64_t readtime;
	int64_t writetime;
	int64_t readlat;
	int64_t writelat;
};

extern int readdrivestat(DriveStats a[], int n);
