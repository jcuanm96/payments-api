package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/utils"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

const errValidatingSendbirdChannelf = "%s is required and must be less than %d characters. Was: %s"

func validateChannelID(paramName string, channelID string) error {
	if utils.IsEmptyValue(channelID) || len(channelID) > constants.MAX_SENDBIRD_CHANNEL_ID_LEN {
		return httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf(errValidatingSendbirdChannelf, paramName, constants.MAX_SENDBIRD_CHANNEL_ID_LEN, channelID),
		)
	}
	return nil
}

func decodeCreateGoatChatChannel(c *fiber.Ctx) (interface{}, error) {
	var p request.CreateGoatChatChannel

	c.BodyParser(&p)

	if utils.IsEmptyValue(p.OtherUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"userID is required",
		)
	}

	return p, nil
}

func decodeStartGoatChat(c *fiber.Ctx) (interface{}, error) {
	var p request.StartGoatChat

	c.BodyParser(&p)

	if err := validateChannelID("sendBirdChannelID", p.SendBirdChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeEndGoatChat(c *fiber.Ctx) (interface{}, error) {
	var p request.EndGoatChat

	c.BodyParser(&p)

	if err := validateChannelID("sendBirdChannelID", p.SendBirdChannelID); err != nil {
		return nil, err
	}

	if utils.IsEmptyValue(p.StartTS) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong ending creator chat.",
			"startTS is required",
		)
	}

	return p, nil
}

func decodeConfirmGoatChat(c *fiber.Ctx) (interface{}, error) {
	var p request.ConfirmGoatChat

	c.BodyParser(&p)

	if err := validateChannelID("sendBirdChannelID", p.SendBirdChannelID); err != nil {
		return nil, err
	}

	if utils.IsEmptyValue(p.CustomerUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong confirming creator chat.",
			"customerUserID is required",
		)
	}

	return p, nil
}

func decodeGetGoatChatPostMessages(c *fiber.Ctx) (interface{}, error) {
	var p request.GetGoatChatPostMessages

	c.QueryParser(&p)

	if utils.IsEmptyValue(p.PostID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong getting post messages.",
			"postID is required",
		)
	}

	if p.Limit == 0 {
		p.Limit = 50
	}

	return p, nil
}

func decodeGetChannelWithUser(c *fiber.Ctx) (interface{}, error) {
	var p request.GetChannelWithUser

	c.QueryParser(&p)

	if utils.IsEmptyValue(p.UserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong getting channel with user.",
			"userID is required",
		)
	}

	return p, nil
}

func decodeCreatePaidGroupChannel(c *fiber.Ctx) (interface{}, error) {
	var p request.CreatePaidGroupChannel

	c.BodyParser(&p)

	p.Name = strings.TrimSpace(strings.Replace(p.Name, "\n", " ", -1))
	p.LinkSuffix = strings.TrimSpace(p.LinkSuffix)

	if p.LinkSuffix != "" {
		regex := regexp.MustCompile(constants.LINKSUFFIX_REGEX)
		regexMatch := regex.MatchString(p.LinkSuffix)
		if !regexMatch {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				fmt.Sprintf(constants.LINKSUFFIX_REGEX_ERR, "Link path"),
				fmt.Sprintf("linkSuffix %s did not match regex %s", p.LinkSuffix, constants.LINKSUFFIX_REGEX),
			)
		}
	}

	if p.Currency != constants.DEFAULT_CURRENCY {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"We only support US dollars right now.",
			"We only support US dollars right now.",
		)
	} else if p.PriceInSmallestDenom < constants.MIN_PAID_GROUP_CHAT_PRICE_USD_IN_SMALLEST_DENOM ||
		p.PriceInSmallestDenom > constants.MAX_PAID_GROUP_CHAT_PRICE_USD_IN_SMALLEST_DENOM {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please input a valid price. $3.00 - $1,000.00",
			"Please input a valid price. $3.00 - $1,000.00",
		)
	}

	file, formFileErr := c.FormFile("coverFile")
	if formFileErr != nil && p.CoverFile != nil {
		return nil, httperr.New(
			500,
			http.StatusInternalServerError,
			"Something went wrong when creating paid group chat",
			fmt.Sprintf("Error parsing cover file: %v", formFileErr),
		)
	} else if file != nil {
		p.CoverFile = file
	}

	if len(p.Name) > constants.MAX_GROUP_NAME_LENGTH {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"The name of the channel is too long.",
			fmt.Sprintf("The name of the channel is too long, max: %d", constants.MAX_GROUP_NAME_LENGTH),
		)
	} else if utils.IsEmptyValue(p.Name) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"The name of the channel cannot be empty.",
			"The name of the channel cannot be empty.",
		)
	}

	if p.Description != nil && len(*p.Description) > constants.MAX_GROUP_DESCRIPTION_LENGTH {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Description must be less than %d characters", constants.MAX_GROUP_DESCRIPTION_LENGTH),
			fmt.Sprintf("Description must be less than %d characters", constants.MAX_GROUP_DESCRIPTION_LENGTH),
		)
	}

	for i, benefit := range []*string{p.Benefit1, p.Benefit2, p.Benefit3} {
		if benefit != nil && len(*benefit) > constants.MAX_GROUP_BENEFIT_LENGTH {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				fmt.Sprintf("Benefit%d must be less than %d characters", i, constants.MAX_GROUP_BENEFIT_LENGTH),
				fmt.Sprintf("Benefit%d must be less than %d characters", i, constants.MAX_GROUP_BENEFIT_LENGTH),
			)
		}
	}

	if p.MemberLimit != nil && (*(p.MemberLimit) < constants.MIN_GROUP_MEMBER_LIMIT || *(p.MemberLimit) > constants.MAX_GROUP_MEMBER_LIMIT) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"You can only set a member limit of 2 to 2000.",
			"You can only set a member limit of 2 to 2000.",
		)
	}

	if p.MemberLimit == nil {
		defaultMemberLimit := constants.DEFAULT_PAID_GROUP_MEMBER_LIMIT
		p.MemberLimit = &defaultMemberLimit
	}

	if p.IsMemberLimitEnabled == nil {
		defaultIsLimitEnabled := constants.DEFAULT_PAID_GROUP_IS_MEMBER_LIMIT_ENABLED
		p.IsMemberLimitEnabled = &defaultIsLimitEnabled
	}

	if len(p.Members) > *p.MemberLimit {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Too many members in group.",
			"length of members exceeded memberLimit.",
		)
	}

	return p, nil
}

func decodeJoinPaidGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.JoinPaidGroup
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeUpdatePaidGroupMemberLimits(c *fiber.Ctx) (interface{}, error) {
	var p request.UpdatePaidGroupMemberLimits
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	if p.IsMemberLimitEnabled == nil && p.MemberLimit == nil {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"The member limit and the member limit toggle can't both be empty.",
			"isMemberLimitEnabled and memberLimit can't both be nil.",
		)
	}

	if p.MemberLimit != nil && (*(p.MemberLimit) < constants.MIN_GROUP_MEMBER_LIMIT || *(p.MemberLimit) > constants.MAX_GROUP_MEMBER_LIMIT) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"You can only set a member limit of 2 to 2000.",
			"You can only set a member limit of 2 to 2000.",
		)
	}

	return p, nil
}

func decodeUpdatePaidGroupPrice(c *fiber.Ctx) (interface{}, error) {
	var p request.UpdatePaidGroupPrice
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	if p.Currency != constants.DEFAULT_CURRENCY {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"We only support US dollars right now.",
			"We only support US dollars right now.",
		)
	} else if p.PriceInSmallestDenom < constants.MIN_PAID_GROUP_CHAT_PRICE_USD_IN_SMALLEST_DENOM ||
		p.PriceInSmallestDenom > constants.MAX_PAID_GROUP_CHAT_PRICE_USD_IN_SMALLEST_DENOM {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Please input a valid price. $3.00 - $1,000.00",
			"Please input a valid price. $3.00 - $1,000.00",
		)
	}

	return p, nil
}

func decodeUpdateGroupMetadata(c *fiber.Ctx) (interface{}, error) {
	var p request.UpdateGroupMetadata
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	if p.Description != nil && len(*p.Description) > constants.MAX_GROUP_DESCRIPTION_LENGTH {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Description must be less than %d characters", constants.MAX_GROUP_DESCRIPTION_LENGTH),
			fmt.Sprintf("Description must be less than %d characters", constants.MAX_GROUP_DESCRIPTION_LENGTH),
		)
	}

	for i, benefit := range []*string{p.Benefit1, p.Benefit2, p.Benefit3} {
		if benefit != nil && len(*benefit) > constants.MAX_GROUP_BENEFIT_LENGTH {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				fmt.Sprintf("Benefit%d must be less than %d characters", i, constants.MAX_GROUP_BENEFIT_LENGTH),
				fmt.Sprintf("Benefit%d must be less than %d characters", i, constants.MAX_GROUP_BENEFIT_LENGTH),
			)
		}
	}

	return p, nil
}

func decodeUpdateGroupLink(c *fiber.Ctx) (interface{}, error) {
	var p request.UpdateGroupLink
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	regex := regexp.MustCompile(constants.LINKSUFFIX_REGEX)
	regexMatch := regex.MatchString(p.LinkSuffix)
	if !regexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf(constants.LINKSUFFIX_REGEX_ERR, "Link path"),
			fmt.Sprintf("linkSuffix %s did not match regex %s", p.LinkSuffix, constants.LINKSUFFIX_REGEX),
		)

	}

	p.LinkSuffix = strings.TrimSpace(p.LinkSuffix)

	return p, nil
}

func decodeCancelPaidGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.CancelPaidGroup
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeScheduledRemoveFromPaidGroup(c *fiber.Ctx) (interface{}, error) {
	var p cloudtasks.RemoveFromPaidGroupTask
	c.BodyParser(&p)

	if v := validate.Struct(p); !v.Validate() {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			v.Errors.One(),
		)
	}

	return p, nil
}

func decodeLeaveGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.LeaveGroup
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeListGoatPaidGroups(c *fiber.Ctx) (interface{}, error) {
	var p request.ListGoatPaidGroups
	c.QueryParser(&p)

	if utils.IsEmptyValue(p.GoatID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong listing creator's paid groups.",
			"goatID is required",
		)
	}
	if p.Limit == 0 {
		p.Limit = 10
	}

	return p, nil
}

func decodeDeletePaidGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.DeletePaidGroup
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeGetPaidGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.GetPaidGroup
	c.QueryParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeBanUserFromPaidGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.BanUserFromPaidGroup
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.BannedUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when banning user from paid group.",
			"userID is required.",
		)
	}

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeUnbanUserFromPaidGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.UnbanUserFromPaidGroup
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.BannedUserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when unbanning user from paid group.",
			"userID is required.",
		)
	}

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeCheckLinkSuffixIsTaken(c *fiber.Ctx) (interface{}, error) {
	var p request.CheckLinkSuffixIsTaken
	c.QueryParser(&p)

	p.LinkSuffix = strings.TrimSpace(p.LinkSuffix)
	regex := regexp.MustCompile(constants.LINKSUFFIX_REGEX)
	regexMatch := regex.MatchString(p.LinkSuffix)
	if !regexMatch {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf(constants.LINKSUFFIX_REGEX_ERR, "Link path"),
			fmt.Sprintf("linkSuffix %s did not match regex %s", p.LinkSuffix, constants.LINKSUFFIX_REGEX),
		)
	}

	return p, nil
}

func decodeListBannedUsers(c *fiber.Ctx) (interface{}, error) {
	var p request.ListBannedUsers
	c.QueryParser(&p)

	if utils.IsEmptyValue(p.ChannelID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong getting banned users",
			"channelID is required.",
		)
	}

	if p.Limit == 0 {
		p.Limit = 10
	}

	return p, nil
}

func decodeRemoveUserFromGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.RemoveUserFromGroup
	c.BodyParser(&p)

	if utils.IsEmptyValue(p.UserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when removing user from group.",
			"userID is required.",
		)
	}

	if utils.IsEmptyValue(p.ChannelID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong when removing user from group.",
			"channelID is required.",
		)
	}

	return p, nil
}

func decodeGetDeepLinkInfo(c *fiber.Ctx) (interface{}, error) {
	var p request.GetDeepLinkInfo
	c.QueryParser(&p)

	if utils.IsEmptyValue(p.LinkSuffix) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Sorry, this user doesn't seem to exist.",
			"linkSuffix is required.",
		)
	}

	return p, nil
}

func decodeGetBatchMediaDownloadURLs(c *fiber.Ctx) (interface{}, error) {
	var q request.GetBatchMediaDownloadURLs
	c.QueryParser(&q)

	if utils.IsEmptyValue(q.ObjectIDs) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"objectIDs is required and non-empty.",
		)
	}

	const maxObjectIDsLen = 200
	if len(q.ObjectIDs) > maxObjectIDsLen {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			fmt.Sprintf("objectIDs max length is %d.", maxObjectIDsLen),
		)
	}

	for i, objectID := range q.ObjectIDs {
		if objectID == "" {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				constants.ErrSomethingWentWrong,
				fmt.Sprintf("objectID at index %d is empty.", i),
			)
		}
	}

	return q, nil
}

func decodeVerifyMediaObject(c *fiber.Ctx) (interface{}, error) {
	var q request.VerifyMediaObject
	c.QueryParser(&q)

	if utils.IsEmptyValue(q.ObjectID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"objectID is required.",
		)
	}

	return q, nil
}

func decodeCreateFreeGroupChannel(c *fiber.Ctx) (interface{}, error) {
	var p request.CreateFreeGroupChannel

	c.BodyParser(&p)

	p.Name = strings.TrimSpace(strings.Replace(p.Name, "\n", " ", -1))
	p.LinkSuffix = strings.TrimSpace(p.LinkSuffix)

	if p.LinkSuffix != "" {
		regex := regexp.MustCompile(constants.LINKSUFFIX_REGEX)
		regexMatch := regex.MatchString(p.LinkSuffix)
		if !regexMatch {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				fmt.Sprintf(constants.LINKSUFFIX_REGEX_ERR, "Link path"),
				fmt.Sprintf("linkSuffix %s did not match regex %s", p.LinkSuffix, constants.LINKSUFFIX_REGEX),
			)
		}
	}

	file, formFileErr := c.FormFile("coverFile")
	if formFileErr != nil && p.CoverFile != nil {
		return nil, httperr.New(
			500,
			http.StatusInternalServerError,
			"Something went wrong when creating free group chat",
			fmt.Sprintf("Error parsing cover file: %v", formFileErr),
		)
	} else if file != nil {
		p.CoverFile = file
	}

	if len(p.Name) > constants.MAX_GROUP_NAME_LENGTH {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("The name of the channel is too long, max: %d", constants.MAX_GROUP_NAME_LENGTH),
			fmt.Sprintf("The name of the channel is too long, max: %d", constants.MAX_GROUP_NAME_LENGTH),
		)
	} else if utils.IsEmptyValue(p.Name) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"The name of the channel cannot be empty.",
			"The name of the channel cannot be empty.",
		)
	}

	if p.Description != nil && len(*p.Description) > constants.MAX_GROUP_DESCRIPTION_LENGTH {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			fmt.Sprintf("Description must be less than %d characters", constants.MAX_GROUP_DESCRIPTION_LENGTH),
			fmt.Sprintf("Description must be less than %d characters", constants.MAX_GROUP_DESCRIPTION_LENGTH),
		)
	}

	for i, benefit := range []*string{p.Benefit1, p.Benefit2, p.Benefit3} {
		if benefit != nil && len(*benefit) > constants.MAX_GROUP_BENEFIT_LENGTH {
			return nil, httperr.New(
				400,
				http.StatusBadRequest,
				fmt.Sprintf("Benefit%d must be less than %d characters", i, constants.MAX_GROUP_BENEFIT_LENGTH),
				fmt.Sprintf("Benefit%d must be less than %d characters", i, constants.MAX_GROUP_BENEFIT_LENGTH),
			)
		}
	}

	if p.MemberLimit != nil && (*(p.MemberLimit) < constants.MIN_GROUP_MEMBER_LIMIT || *(p.MemberLimit) > constants.MAX_GROUP_MEMBER_LIMIT) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"You can only set a member limit of 2 to 2000.",
			"You can only set a member limit of 2 to 2000.",
		)
	}

	if p.MemberLimit == nil {
		defaultMemberLimit := constants.DEFAULT_PAID_GROUP_MEMBER_LIMIT
		p.MemberLimit = &defaultMemberLimit
	}

	if p.IsMemberLimitEnabled == nil {
		defaultIsLimitEnabled := constants.DEFAULT_PAID_GROUP_IS_MEMBER_LIMIT_ENABLED
		p.IsMemberLimitEnabled = &defaultIsLimitEnabled
	}

	if len(p.Members) > *p.MemberLimit {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Too many members in group.",
			"length of members exceeded memberLimit.",
		)
	}

	return p, nil
}

func decodeJoinFreeGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.JoinFreeGroup
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeGetFreeGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.GetFreeGroup
	c.QueryParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeListUserFreeGroups(c *fiber.Ctx) (interface{}, error) {
	var p request.ListUserFreeGroups
	c.QueryParser(&p)

	if utils.IsEmptyValue(p.UserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			"Something went wrong listing user's free groups.",
			"userID is required",
		)
	}
	if p.Limit == 0 {
		p.Limit = 10
	}

	return p, nil
}

func decodeDeleteFreeGroup(c *fiber.Ctx) (interface{}, error) {
	var p request.DeleteFreeGroup
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	return p, nil
}

func decodeAddFreeGroupCoCreators(c *fiber.Ctx) (interface{}, error) {
	var p request.AddFreeGroupChatCoCreators
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	if utils.IsEmptyValue(p.CoCreators) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"coCreators cannot be empty.",
		)
	}

	return p, nil
}

func decodeRemoveFreeGroupCoCreator(c *fiber.Ctx) (interface{}, error) {
	var p request.RemoveFreeGroupChatCoCreator
	c.BodyParser(&p)

	if err := validateChannelID("channelID", p.ChannelID); err != nil {
		return nil, err
	}

	if utils.IsEmptyValue(p.UserID) {
		return nil, httperr.New(
			400,
			http.StatusBadRequest,
			constants.ErrSomethingWentWrong,
			"userID is required.",
		)
	}

	return p, nil
}
