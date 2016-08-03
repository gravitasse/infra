#ifndef ONLP_H
#define ONLP_H

#include "pluginCommon.h"

int Init();
int DeInit();
int GetMaxNumOfFans();
int GetAllFanState(fan_info_t *, int);
int GetFanState(fan_info_t *, int);
int GetPlatformState(sys_info_t *);

#endif //ONPL_H
