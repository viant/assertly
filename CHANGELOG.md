## April 23 2019 - v0.4.7
    * patched time validation with custom time format

## April 23 2019 - v0.4.6
    * Updated keysValue for selector returing a map

## April 2 2019 - v0.4.5
    * Updated multi line JSON handling 
    
## Feb 23 2019 - v0.4.4
    * Added limit to available keys

## Feb 12 2019 - v0.4.3
    * Patched @numericPrecisionPoint@ with zero value

## Feb 10 2019 - v0.4.2
    * Added empty value provider: <ds:empty>

## Jan 30 2019 - v0.4.1
   * Patch assertPath directive
   * Remove cast error login, in case of error original value is used
    
## Nov 30 2018 - v0.3.0
   * Added @length@ directive
   * Patched map level directives

## Nov 30 2018 - v0.3.0
   * Added assertPath directive

## Nov 30 2018 - v0.2.2
   * Update assertSlice logic to account for key/value pairs as map validation
   * Added keyCaseSensitive directive
   * Change behaviour to caseSensitive directive to only apply to values
   * Added coalesceWithZero directive, patched 0 | 0.0 with nil validation

## Nov 26 2018 - v0.2.1
   * Added numericPrecisionPoint directive

## Nov 26 2018 - v0.2.0
   * Added numericPrecisionPoint directive

## Nov 25 2018 - v0.1.1
   * Patch nil pointer issue on assertTime
   * Added int check reporting if applicable in assertFloat
    
## Nov 5 2018 - v0.1.0
  * Enhanced apply time directory to only non function based values
  * Introduced AssertValuesWithContext helper validation method with context
  
## Jan 17 2018 (Alpha)
  * Initial Release.
