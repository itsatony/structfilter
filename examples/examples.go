package main

import (
	"fmt"
	"time"

	"github.com/itsatony/structfilter"
)

type SourceStruct struct {
	Field1 string `filter:"public"`
	Field2 int    `filter:"private,admin"`
	Field3 bool   `filter:"public,user"`
	Field4 bool
}

type QuizQuestion struct {
	Id                                string             `json:"id" gorm:"primaryKey" filter:"public,nospoiler"`
	TenantId                          string             `json:"tenantId" gorm:"secondaryKey" filter:"public,nospoiler"`
	Tags                              []string           `json:"tags" gorm:"type:text[];index" filter:"public,nospoiler"`
	Title                             string             `json:"title" filter:"public,nospoiler"`
	Type                              string             `json:"type" filter:"public,nospoiler"` // SingleChoiceFromCollectionAnswerOptions, SingleChoice, MultipleChoice, FreeText, OneFloat, TwoFloats, OneString, TwoStrings, Binary, Sorting
	SubmissionTimeStart               time.Time          `json:"submissionTimeStart" filter:"public,nospoiler"`
	SubmissionTimeEnd                 time.Time          `json:"submissionTimeEnd" filter:"public,nospoiler"`
	VisibleTimeStart                  time.Time          `json:"visibleTimeStart" filter:"public,nospoiler"`
	VisibleTimeEnd                    time.Time          `json:"visibleTimeEnd" filter:"public,nospoiler"`
	IsFreeTextValidationCaseSensitive bool               `json:"isFreeTextValidationCaseSensitive" filter:"public,nospoiler"`
	IsPublic                          bool               `json:"isPublic" gorm:"index" filter:"public,nospoiler"`
	Description                       string             `json:"description" filter:"public,nospoiler"`
	Hint                              string             `json:"hint" filter:"public,nospoiler"`
	Explanation                       string             `json:"explanation" filter:"public,nospoiler"`
	QuestionText                      string             `json:"questionText" filter:"public,nospoiler"`
	QuestionMediaUrl                  string             `json:"mediaUrl" filter:"public,nospoiler"`
	QuestionText2                     string             `json:"questionText2" filter:"public,nospoiler"`
	QuestionMediaUrl2                 string             `json:"mediaUrl2" filter:"public,nospoiler"`
	CreatedAt                         time.Time          `json:"createdAt" filter:"public,nospoiler"`
	CreatedByUserId                   string             `json:"createdByUserId,omitempty" filter:"admin"`
	UpdatedAt                         time.Time          `json:"updatedAt" filter:"public,nospoiler"`
	UpdatedByUserId                   string             `json:"updatedByUserId,omitempty" filter:"admin"`
	DeletedAt                         time.Time          `json:"-" gorm:"index" filter:"admin"`
	IsPrediction                      bool               `json:"isPrediction" filter:"public,nospoiler"`
	IsResolved                        bool               `json:"isResolved" gorm:"index" filter:"public,nospoiler"`
	SolutionAnswerIds                 []string           `json:"solutionAnswerIds" gorm:"type:text[];index" filter:"public,spoiler"`
	SolutionAnswerString01            string             `json:"solutionAnswerString01" gorm:"index" filter:"public,spoiler"`
	SolutionAnswerString02            string             `json:"solutionAnswerString02" gorm:"index" filter:"public,spoiler"`
	SolutionAnswerFloat01             float64            `json:"solutionAnswerFloat01" gorm:"index" filter:"public,spoiler"`
	SolutionAnswerFloat02             float64            `json:"solutionAnswerFloat02" gorm:"index" filter:"public,spoiler"`
	AnswerOptions                     []QuizAnswerOption `json:"answerOptions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:question_id;references:id" filter:"public,nospoiler"`
	GlobalQuestionId                  string             `json:"globalQuestionId" filter:"public,nospoiler"`
} //@name QuizQuestion

type QuizAnswerOption struct {
	Id                   string    `json:"id" gorm:"primaryKey" filter:"public,nospoiler"`
	TenantId             string    `json:"tenantId" gorm:"secondaryKey" filter:"public,nospoiler"`
	AnswerText           string    `json:"answerText" filter:"public,nospoiler"`
	AnswerMediaUrl       string    `json:"mediaUrl" filter:"public,nospoiler"`
	IsCorrect            bool      `json:"isCorrect" filter:"public,spoiler"`
	CreatedAt            time.Time `json:"createdAt" filter:"public,nospoiler"`
	CreatedByUserId      string    `json:"createdByUserId,omitempty" filter:"admin"`
	UpdatedAt            time.Time `json:"updatedAt" filter:"public,nospoiler"`
	UpdatedByUserId      string    `json:"updatedByUserId,omitempty" filter:"admin"`
	DeletedAt            time.Time `json:"deletedAt" gorm:"index" filter:"admin"`
	GlobalAnswerOptionId string    `json:"globalAnswerOptionId" filter:"public,nospoiler"`
	QuestionId           string    `json:"questionId" gorm:"index"`
	CollectionId         string    `json:"collectionId" gorm:"index"` // this is for cases where each question in a collection should have the same answer-option-pool
} //@name QuizAnswerOption

func main() {
	// filter to keep all public fields
	source := SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered := structfilter.CreateFilteredStruct(source, []string{"public"}, nil)
	fmt.Println("(1) filter to keep all public fields:", filtered)
	// filter to remove all admin fields
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered = structfilter.CreateFilteredStruct(source, []string{""}, []string{"admin"})
	fmt.Println("(2) filter to remove all admin fields and keep any others:", filtered)
	// filter to keep all admin fields
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered = structfilter.CreateFilteredStruct(source, []string{"admin"}, nil)
	fmt.Println("(3) filter to keep all admin fields:", filtered)
	// filter to keep all fields with any filter value
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered = structfilter.CreateFilteredStruct(source, []string{""}, nil)
	fmt.Println("(4) filter to keep all fields with any filter value:", filtered)
	// EmptyFields example
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: true,
	}
	emptyFields, err := structfilter.EmptyFilteredFields(&source, map[string][]string{"filter": {"admin"}})
	if err != nil {
		panic(err)
	}
	fmt.Println("(5) filter to empty the fields listed BEFORE:", source)
	fmt.Println("(5) filter to empty the fields listed AFTER:", emptyFields)

	question := QuizQuestion{
		Id:                                "1",
		TenantId:                          "1",
		SubmissionTimeStart:               time.Now(),
		SubmissionTimeEnd:                 time.Now(),
		VisibleTimeStart:                  time.Now(),
		VisibleTimeEnd:                    time.Now(),
		IsPublic:                          true,
		IsFreeTextValidationCaseSensitive: true,
		Description:                       "desc public",
		Hint:                              "hint public",
		Explanation:                       "explanation public",
		QuestionText:                      "qtext",
		QuestionMediaUrl:                  "qmediaurl",
		QuestionText2:                     "asasfd",
		QuestionMediaUrl2:                 "asasfd",
		CreatedAt:                         time.Now(),
		CreatedByUserId:                   "1",
		UpdatedAt:                         time.Now(),
		UpdatedByUserId:                   "1",
		DeletedAt:                         time.Now(),
		IsPrediction:                      true,
		IsResolved:                        true,
		SolutionAnswerIds:                 []string{"1", "2"},
		SolutionAnswerString01:            "sasas",
		SolutionAnswerString02:            "sasas",
		SolutionAnswerFloat01:             1.0,
		SolutionAnswerFloat02:             1.0,
		GlobalQuestionId:                  "1",
		AnswerOptions: []QuizAnswerOption{
			{
				Id:                   "1",
				TenantId:             "1",
				AnswerText:           "answertext",
				AnswerMediaUrl:       "answertext",
				IsCorrect:            true,
				CreatedAt:            time.Now(),
				CreatedByUserId:      "1",
				UpdatedAt:            time.Now(),
				UpdatedByUserId:      "1",
				DeletedAt:            time.Now(),
				GlobalAnswerOptionId: "1",
				QuestionId:           "1",
				CollectionId:         "1",
			},
		},
	}
	emptyfields, err := structfilter.EmptyFilteredFields(&question, map[string][]string{"filter": {"spoiler"}})
	if err != nil {
		panic(err)
	}
	filteredQuestion, castOkay := emptyfields.(*QuizQuestion)
	if !castOkay {
		fmt.Printf("Could not cast filtered question to QuizQuestion\n")
		return
	}
	fmt.Printf("Filtered question: %v\n", filteredQuestion)
}
