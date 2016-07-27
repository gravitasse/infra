#ifndef ONLP_H
#define ONLP_H

#include "pluginCommon.h"

int Init();
int GetMaxNumOfFans();
int GetAllFanState(fan_info_t *info, int count);
int GetFanState(fan_info_t *info, int fanId);
#endif //ONPL_H
