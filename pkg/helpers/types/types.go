package types

type PathTypes int

const PathTypesCount = 13
const (
	_ PathTypes = iota
	UUIDPath
	NounPath
	TwoNounPath
	ThreeNounPath
	TwoNounDashedPath
	ThreeNounDashedPath
	VerbPath
	VerbNounPath
	DatePath
	YearPath
	MonthPath
	DayPath
)

type VariableNames int

const VariableNamesCount = 12
const (
	_ VariableNames = iota
	SingleCharacterVariable
	VerbVariable
	NounVariable
	TwoNounVariable
	ThreeNounVariable
	TwoNounDashedVariable
	ThreeNounDashedVariable
	VerbNounCombinationVariable
	NounVerbCombinationVariable
	MonthVariable
	DayVariable
)

type VariableValues int

const VariableValuesCount = 14
const (
	_ = iota
	VerbValue
	NounValue
	TwoNounValue
	ThreeNounValue
	TwoNounDashedValue
	ThreeNounDashedValue
	VerbNounCombinationValue
	NounVerbCombinationValue
	DateValue
	YearValue
	MonthValue
	DayValue
	Base64Value
)
