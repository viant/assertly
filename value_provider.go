package assertly

import (
	"github.com/viant/toolbox"
)

//ValueProviderRegistry represents value provider ValueProviderRegistry
var ValueProviderRegistry = toolbox.NewValueProviderRegistry()

func init() {
	ValueProviderRegistry.Register("nil", toolbox.NewNilValueProvider())
	ValueProviderRegistry.Register("env", toolbox.NewEnvValueProvider())
	ValueProviderRegistry.Register("cast", toolbox.NewCastedValueProvider())
	ValueProviderRegistry.Register("timediff", toolbox.NewTimeDiffProvider())
	ValueProviderRegistry.Register("current_timestamp", toolbox.NewCurrentTimeProvider())
	ValueProviderRegistry.Register("current_date", toolbox.NewCurrentDateProvider())
	ValueProviderRegistry.Register("between", toolbox.NewBetweenPredicateValueProvider())
	ValueProviderRegistry.Register("within_sec", toolbox.NewWithinSecPredicateValueProvider())
	ValueProviderRegistry.Register("weekday", toolbox.NewWeekdayProvider())
	ValueProviderRegistry.Register("dob", toolbox.NewDateOfBirthrovider())
	ValueProviderRegistry.Register("cat", toolbox.NewFileValueProvider(true))

}
